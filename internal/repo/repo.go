package repo

import (
	"github.com/oanhngy/url-shortener/internal/model"
)

type LinkRepository interface {
	Save(link *model.Link) error //lưu vào db
	FindByCode(code string) (*model.Link, error)
	Exists(code string) bool //check collision
	FindAll() ([]model.Link, error)
	IncrementClick(code string) error //click count
}
