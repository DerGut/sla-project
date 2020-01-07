package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const redisHost = "cache"
const redisPort = "6379"
const redisKey = "featured"

type Cache interface {
	Close()
	GetFeaturedData(db DB) (*[]Document, error)
	VoteUp(primitive.ObjectID) (int, error)
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() (r *redisCache, err error) {
	client := redis.NewClient(&redis.Options{Addr: redisHost + ":" + redisPort})

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	log.Printf("Initialized connection to redis")
	return &redisCache{client}, nil
}

func (r *redisCache) Close() {
	r.client.Close()
}

func (r *redisCache) getFeaturedIds() (ids []primitive.ObjectID, scores []int, err error) {
	opt := redis.ZRangeBy{
		Min:    "0",
		Max:    "+inf",
		Offset: 0,
		Count:  5,
	}
	featured, err := r.client.ZRevRangeByScoreWithScores(redisKey, opt).Result()
	if err != nil {
		return nil, nil, err
	}

	for _, x := range featured  {
		oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", x.Member))
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, oid)
		scores = append(scores, int(x.Score))
	}
	
	return ids, scores, nil
}

func (r *redisCache) updateFromDB(db DB, id primitive.ObjectID) (doc *Document, err error) {
	doc, err = db.FindOne(id)
	if err != nil {
		return nil, err
	}
	bin, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	_, err = r.client.Set(id.Hex(), string(bin), 1*time.Hour).Result()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (r *redisCache) GetFeaturedData(db DB) (docs *[]Document, err error) {
	// get featured ids with scores
	ids, scores, err := r.getFeaturedIds()
	if err != nil {
		return nil, err
	}
	
	// get featured docs from redis
	var d []Document
	for i, id := range ids {
		val, err := r.client.Get(id.Hex()).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		if err == redis.Nil {
			log.Printf("Document is not cached in redis - retrieving from mongo: %s", id.Hex())
			doc, err := r.updateFromDB(db, id)
			if err != nil {
				return nil, err
			}
			d = append(d, *doc)
		} else {
			log.Printf("Using cached document from redis")
			var doc Document
			err = json.Unmarshal([]byte(val), &doc)
			if err != nil {
				log.Printf("Something went wrong doing unmarshalling")
				return nil, err
			}
			doc.Id = id
			doc.Upvotes = scores[i]
			d = append(d, doc)
		}
	}
	
	return &d, nil
}

func (r *redisCache) VoteUp(id primitive.ObjectID) (int, error) {
	upvotes, err := r.client.ZIncrBy(redisKey, 1, id.Hex()).Result()
	if err != nil {
		return -1, err
	}

	return int(upvotes), nil
}
