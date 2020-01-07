package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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

	t, err = template.ParseFiles("templates/index.html", "templates/table.html")
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
		log.Printf("Couoldn't read request body: %s", err)
		return
	}

	id, err := primitive.ObjectIDFromHex(string(body))
	if err != nil {
		log.Printf("No valid ObjectId: %s", err)
		http.NotFound(w, r)
		return
	}

	upvotes, err := cache.VoteUp(id)
	if err != nil {
		log.Printf("Couldn't vote up in redis: %s", err)
	}
	w.Write([]byte(strconv.Itoa(upvotes)))

	err = db.VoteUp(id)
	if err != nil {
		log.Printf("Couldn't vote up in mongo: %s", err)
	}
}

func featuredDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("featuredDataHandler handling %s", r.URL.Path)

	featured, err := cache.GetFeaturedData(db)
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
		log.Printf("Couldn't marshal response: %s", err)
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
		log.Printf("Couldn't marshal response: %s", err)
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

	featured, err := cache.GetFeaturedData(db)
	if err != nil {
		log.Printf("Couldn't find featured data in redis - trying mongo: %s", err)
		featured, err = db.FindFeaturedData()
		if err != nil {
			log.Printf("Couldn't find featured data in mongo either: %s", err)
			featured = &[]Document{}
		}
	}

	all, err := db.FindImportantData()
	if err != nil {
		log.Printf("Couldn't find important data: %s", err)
		all = &[]Document{}
	}

	d := TemplateData{
		Title: "Sample App",
		Featured: *featured,
		All:  *all,
	}
	t.Execute(w, d)
}

func main() {
	defer db.Close()
	defer cache.Close()
	defer queue.Close()

	err := queue.PublishTask(Document{Val1:"rwe", Val2:"dasd"})
	if err != nil {
		log.Fatalf("Couldn't publish: %s", err)
	}

	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/upvote/", upvoteHandler)
	http.HandleFunc("/featured-data/", featuredDataHandler)
	http.HandleFunc("/all-data/", allDataHandler)
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

