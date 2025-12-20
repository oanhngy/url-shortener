package main

import (
	"log"
	"net/http"

	"github.com/oanhngy/url-shortener/internal/handler"
	"github.com/oanhngy/url-shortener/internal/repo/memory"
	"github.com/oanhngy/url-shortener/internal/service"
)

func main() {
	repo := memory.NewInMemoryRepo()    //repo
	svc := service.NewLinkService(repo) //service
	baseURL := "http://localhost:8080"
	h := handler.NewLinkHandler(svc, baseURL) //handler
	mux := http.NewServeMux()                 //tạo router

	//endpoint tạo link +link list
	mux.HandleFunc("/api/v1/links", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost { //POST
			h.CreateLink(w, r)
			return
		}
		if r.Method == http.MethodGet { //GET
			h.ListLinks(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) //405
	})

	//endpoint view info
	mux.HandleFunc("/api/v1/links/", h.ViewInfo) //GET

	//endpoint redirect
	mux.HandleFunc("/", h.Redirect) //GET

	addr := ":8080"
	log.Printf("Starting server at %s\n", addr)

	//start server
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
