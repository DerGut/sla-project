package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

type Document struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`

	Val1    string             `bson:"val1" json:"val1"`
	Val2    string             `bson:"val2" json:"val2"`
	Val3    string             `bson:"val3" json:"val3"`

	Upvotes int                `bson:"upvotes" json:"upvotes"`
}

type TemplateData struct {
	Title 	 	  string
	Featured, All []Document
}

var (
	t     *template.Template
	db    DB
	cache Cache
	queue Queue
)

func init() {
	var err error

	tracer.Start(tracer.WithServiceName("frontend"))

	t, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Couldn't parse templates: %s", err)
	}

	db, err = NewMongoDB()
	if err != nil {
		log.Fatalf("Couldn't connect to mongo: %s", err)
	}

	cache, err = NewRedisCache()
	if err != nil {
		log.Fatalf("Couldn't connect to redis: %s", err)
	}

	queue, err = NewRabbitMQQueue()
	if err != nil {
		log.Fatalf("Couldn't connect to rabbit: %s", err)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("staticHandler handling %s", r.URL.Path)

	path := r.URL.Path[1:]
	content, err := ioutil.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("Couldn't read file %s", path)
		return
	}

	ext := filepath.Ext(path)
	contentTypes := map[string]string{
		".css": "text/css",
		".js":  "text/js",
	}
	if ct, ok := contentTypes[ext]; ok {
		w.Header().Add("Content-Type", ct)
	}
	w.Write(content)
}

func upvoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("upvoteHandler handling %s", r.URL.Path)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %s", err)
	}

	id, err := primitive.ObjectIDFromHex(string(body))
	if err != nil {
		log.Printf("No valid ObjectId %s: %s", string(body), err)
		http.NotFound(w, r)
		return
	}

	err = db.VoteUp(id)
	if err != nil {
		log.Fatalf("Couldn't vote up in mongo: %s", err)
	}
}

func featuredDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("featuredDataHandler handling %s", r.URL.Path)

	featured, err := cache.GetFeaturedData()
	if err != nil {
		log.Printf("Couldn't find featured data in redis - trying mongo: %s", err)
		featured, err = db.FindFeaturedData()
		if err != nil {
			log.Printf("Couldn't find featured data in mongo either: %s", err)
			featured = &[]Document{}
		}
	}

	marshalled, err := json.Marshal(featured)
	if err != nil {
		log.Fatalf("Couldn't marshal response: %s", err)
	}

	w.Write(marshalled)
	w.Header().Add("Content-Type", "text/json")
}

func allDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("allDataHandler handling %s", r.URL.Path)

	all, err := db.FindImportantData()
	if err != nil {
		log.Printf("Couldn't find important data: %s", err)
		all = &[]Document{}
	}

	marshalled, err := json.Marshal(all)
	if err != nil {
		log.Fatalf("Couldn't marshal response: %s", err)
	}

	w.Write(marshalled)
	w.Header().Add("Content-Type", "text/json")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("indexHandler handling %s", r.URL.Path)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	d := TemplateData{
		Title: "Sample App",
	}
	t.Execute(w, d)
}

func main() {
	defer tracer.Stop()
	defer db.Close()
	defer cache.Close()
	defer queue.Close()

	//go keepPublishingTasks()

	go syncCache()

	mux := httptrace.NewServeMux()
	mux.HandleFunc("/static/", staticHandler)
	mux.HandleFunc("/upvote/", upvoteHandler)
	mux.HandleFunc("/featured-data/", featuredDataHandler)
	mux.HandleFunc("/all-data/", allDataHandler)
	mux.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func keepPublishingTasks() {
	for {
		token := make([]byte, 10)
		rand.Read(token)
		err := queue.PublishTask(Document{
			Val1: string(token[:5]),
			Val2: string(token[5:]),
		})
		if err != nil {
			log.Printf("Couldn't publish: %s", err)
		}
		time.Sleep(3*time.Second)
	}
}

func syncCache() {
	for range time.NewTicker(5 * time.Second).C {
		err := syncCacheWithDB(cache, db)
		if err != nil {
			log.Printf("Problem syncing cache: %s", err)
		} else {
			log.Printf("Synced cache")
		}
	}
}

func syncCacheWithDB(c Cache, db DB) error {
	featured, err := db.FindFeaturedData()
	if err != nil {
		return err
	}

	return c.UpdateFeaturedData(featured)
}
