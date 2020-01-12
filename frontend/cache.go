package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
	"log"
)

const redisHost = "cache"
const redisPort = "6379"
const redisFeaturedKey = "featured"

type Cache interface {
	Close()
	GetFeaturedData() (*[]Document, error)
	UpdateFeaturedData(*[]Document) error
}

type redisCache struct {
	client *redistrace.Client
}

func NewRedisCache() (r *redisCache, err error) {
	client := redistrace.NewClient(&redis.Options{Addr: redisHost + ":" + redisPort})

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

func (r *redisCache) GetFeaturedData() (*[]Document, error)  {
	opt := redis.ZRangeBy{
		Min: "0",
		Max:"+inf",
		Offset: 0,
		Count: 5,
	}
	featured, err := r.client.ZRevRangeByScoreWithScores(redisFeaturedKey, opt).Result()
	if err != nil {
		return nil, err
	}

	var docs []Document
	for _, x := range featured {
		var doc Document
		err := json.Unmarshal([]byte(fmt.Sprintf("%v", x.Member)), &doc)
		if err != nil {
			return nil, err
		}
		doc.Upvotes = int(x.Score)
		docs = append(docs, doc)
	}

	return &docs, nil
}

func (r *redisCache) UpdateFeaturedData(docs *[]Document) error {
	var members []redis.Z
	for _, d := range *docs {
		// Neutralize all upvotes so that duplicates will be overwritten
		score := float64(d.Upvotes)
		d.Upvotes = 0

		bytes, err := json.Marshal(d)
		if err != nil {
			return err
		}
		members = append(members, redis.Z{
			Score:  score,
			Member: string(bytes),
		})
	}
	_, err := r.client.ZAdd(redisFeaturedKey, members...).Result()
	return err
}
