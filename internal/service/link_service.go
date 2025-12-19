package service

import (
	"errors"
	"time"

	"github.com/oanhngy/url-shortener/internal/model"
	"github.com/oanhngy/url-shortener/internal/repo"
	"github.com/oanhngy/url-shortener/packages/base62"
)

// xử lý link
type LinkService struct {
	repo repo.LinkRepository
}

// tạo link mới
func NewLinkService(r repo.LinkRepository) *LinkService {
	return &LinkService{repo: r}
}

// tạo shortURL
func (s *LinkService) ShortenURL(longURL string) (model.Link, error) {
	if longURL == "" {
		return model.Link{}, errors.New("long URL is empty")
	}

	var code string
	for {
		code = base62.Generate(6)
		if !s.repo.Exists(code) {
			break
		}
	}

	link := model.Link{
		ID:        "", //tự động tạo ID
		LongURL:   longURL,
		ShortCode: code,
		CreatedAt: time.Now(),
	}

	err := s.repo.Save(&link)
	if err != nil {
		return model.Link{}, err
	}

	return link, nil
}

// trả về longURL từ shortCode
func (s *LinkService) ReturnOriginalURL(code string) (string, error) {
	link, err := s.repo.FindByCode(code)
	if err != nil {
		return "", err
	}

	_ = s.repo.IncrementClick(code) //tăng click count
	return link.LongURL, nil
}

// info's link
func (s *LinkService) ViewInfo(code string) (*model.Link, error) {
	return s.repo.FindByCode(code)
}

// view all links
func (s *LinkService) ViewAllLinks() ([]model.Link, error) {
	return s.repo.FindAll()
}
