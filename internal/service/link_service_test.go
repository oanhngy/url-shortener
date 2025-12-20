package service

import (
	"testing"

	"github.com/oanhngy/url-shortener/internal/repo/memory"
)

func TestShortenURL_Success(t *testing.T) {
	repo := memory.NewInMemoryRepo()
	service := NewLinkService(repo)

	longURL := "https://example.com/abcdef987"

	link, err := service.ShortenURL(longURL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link.LongURL != longURL {
		t.Errorf("expected longURL %s, got %s", longURL, link.LongURL)
	}

	if link.ShortCode == "" {
		t.Errorf("expected shortCode to be generated")
	}
}

func TestShotenURL_NoCollision(t *testing.T) {
	repo := memory.NewInMemoryRepo()
	service := NewLinkService(repo)

	link1, _ := service.ShortenURL("https://a.com")
	link2, _ := service.ShortenURL("https://b.com")

	if link1.ShortCode == link2.ShortCode {
		t.Errorf("expected different short codes, got same %s", link1.ShortCode)
	}
}

func TestReturnOriginalURL_IncreaseClick(t *testing.T) {
	repo := memory.NewInMemoryRepo()
	service := NewLinkService(repo)

	link, _ := service.ShortenURL("https://example.com")
	_, _ = service.ReturnOriginalURL(link.ShortCode)
	_, _ = service.ReturnOriginalURL(link.ShortCode)

	stored, _ := repo.FindByCode(link.ShortCode)

	if stored.ClickCount != 2 {
		t.Errorf("expected click count 2, got %d", stored.ClickCount)
	}
}
