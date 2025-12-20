package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/oanhngy/url-shortener/internal/model"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	links map[string]*model.Link //key=ShortCode, value=Link
}

// tạo repo mới
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		links: make(map[string]*model.Link), //map rỗng
	}
}

// lưu link vào map
func (r *InMemoryRepo) Save(link *model.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if link == nil {
		return errors.New("link is nil")
	}

	if link.ShortCode == "" {
		return errors.New("short code is empty")
	}

	//collision
	if _, ok := r.links[link.ShortCode]; ok {
		return errors.New("short code already exists")
	}

	if link.ID == "" {
		link.ID = link.ShortCode
	}

	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now()
	}

	cp := *link
	r.links[link.ShortCode] = &cp

	return nil
}

// tìm link theo shortCode
func (r *InMemoryRepo) FindByCode(code string) (*model.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.links[code]
	if !ok {
		return nil, errors.New("link not found")
	}

	cp := *link
	return &cp, nil
}

// check shortCode tồn tại chưa, collision
func (r *InMemoryRepo) Exists(code string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.links[code]
	return ok
}

// trả về all links
func (r *InMemoryRepo) FindAll() ([]model.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Link, 0, len(r.links))
	for _, link := range r.links {
		out = append(out, *link)
	}

	//sort theo ngày giảm dần
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// tăng click count
func (r *InMemoryRepo) IncrementClick(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	link, ok := r.links[code]
	if !ok {
		return errors.New("link not found")
	}

	link.ClickCount++
	return nil
}
