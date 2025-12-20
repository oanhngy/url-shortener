package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/oanhngy/url-shortener/internal/service"
)

// chứa service xử lý
type LinkHandler struct {
	svc     *service.LinkService
	baseURL string
}

// handler mới
func NewLinkHandler(svc *service.LinkService, baseURL string) *LinkHandler {
	return &LinkHandler{
		svc:     svc,
		baseURL: baseURL,
	}
}

// **CREATE LINK
// xu ly POST, shorten
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type reqBody struct {
		LongURL string `json:"longUrl"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) //400
		return
	}

	//gọi service shorten
	link, err := h.svc.ShortenURL(req.LongURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //400
		return
	}

	type respBody struct {
		ID         string `json:"id"`
		LongURL    string `json:"longUrl"`
		ShortCode  string `json:"shortCode"`
		ShortURL   string `json:"shortUrl"`
		ClickCount int    `json:"clicks"`
		CreatedAt  string `json:"createdAt"`
	}

	resp := respBody{
		ID:         link.ID,
		LongURL:    link.LongURL,
		ShortCode:  link.ShortCode,
		ShortURL:   strings.TrimRight(h.baseURL, "/") + "/" + link.ShortCode,
		ClickCount: link.ClickCount,
		CreatedAt:  link.CreatedAt.Format(time.RFC3339),
	}

	writeJSON(w, http.StatusCreated, resp) //201
}

// **LIST LINKS
func (h *LinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
		return
	}

	links, err := h.svc.ViewAllLinks()
	if err != nil {
		http.Error(w, "Failed to list links", http.StatusInternalServerError)
		return
	}

	//item cho từng link
	type item struct {
		ID         string `json:"id"`
		LongURL    string `json:"longUrl"`
		ShortCode  string `json:"shortCode"`
		ShortURL   string `json:"shortUrl"`
		ClickCount int    `json:"clicks"`
		CreatedAt  string `json:"createdAt"`
	}

	out := make([]item, 0, len(links)) //slice
	for _, l := range links {
		out = append(out, item{
			ID:         l.ID,
			LongURL:    l.LongURL,
			ShortCode:  l.ShortCode,
			ShortURL:   strings.TrimRight(h.baseURL, "/") + "/" + l.ShortCode,
			ClickCount: l.ClickCount,
			CreatedAt:  l.CreatedAt.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, out) //200
}

// **VIEW INFO
func (h *LinkHandler) ViewInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//lấy code từ /api/v1/links/{code}
	code := strings.TrimPrefix(r.URL.Path, "/api/v1/links/") //tách
	if code == "" || strings.Contains(code, "/") {
		http.NotFound(w, r) //404
		return
	}

	//gọi service
	link, err := h.svc.ViewInfo(code)
	if err != nil {
		http.NotFound(w, r) //404
		return
	}

	type respBody struct {
		ID         string `json:"id"`
		LongURL    string `json:"longUrl"`
		ShortCode  string `json:"shortCode"`
		ShortURL   string `json:"shortUrl"`
		ClickCount int    `json:"clicks"`
		CreatedAt  string `json:"createdAt"`
	}

	resp := respBody{
		ID:         link.ID,
		LongURL:    link.LongURL,
		ShortCode:  link.ShortCode,
		ShortURL:   strings.TrimRight(h.baseURL, "/") + "/" + link.ShortCode,
		ClickCount: link.ClickCount,
		CreatedAt:  link.CreatedAt.Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, resp)
}

// **REDIRECT
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//lấy code từ /{code}
	code := strings.TrimPrefix(r.URL.Path, "/") //tách
	if code == "" {
		http.NotFound(w, r)
		return
	}

	if strings.HasPrefix(code, "api") {
		http.NotFound(w, r)
		return
	}

	//tìm url+tăng click
	longURL, err := h.svc.ReturnOriginalURL(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound) //302

}

// helper function
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
