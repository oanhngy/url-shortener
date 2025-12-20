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

// xu ly POST, shorten
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type reqBody struct {
		LongURL string `json:"long_url"`
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

// helper function
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
