package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oanhngy/url-shortener/internal/handler"
	mongorepo "github.com/oanhngy/url-shortener/internal/repo/mongo"
	"github.com/oanhngy/url-shortener/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:rootpassword@localhost:27017/?authSource=admin"
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "urlshortener"
	}

	//connect mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	repo := mongorepo.NewMongoRepo(db)

	//unique index=short_code
	if err := repo.EnsureIndexes(); err != nil {
		log.Fatal(err)
	}

	svc := service.NewLinkService(repo)
	baseURL := "http://localhost:8080"
	h := handler.NewLinkHandler(svc, baseURL)

	mux := http.NewServeMux()

	//endpoint tạo link + list links
	mux.HandleFunc("/api/v1/links", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateLink(w, r)
			return
		}
		if r.Method == http.MethodGet {
			h.ListLinks(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	//endpoint view info
	mux.HandleFunc("/api/v1/links/", h.ViewInfo)

	//endpoint redirect
	mux.HandleFunc("/", h.Redirect)

	addr := ":8080"
	log.Printf("Starting server at %s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
