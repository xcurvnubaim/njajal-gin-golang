package shortlink

import (
	"log"

	"gorm.io/gorm"
)

type IRepository interface {
	CreateShortenerLink(data *ShortenerLinkModel) error
	GetShortenerLinkByShortenerURL(shortenerURL string) (*ShortenerLinkModel, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateShortenerLink(data *ShortenerLinkModel) error {
	err := r.db.Create(data).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *repository) GetShortenerLinkByShortenerURL(shortenerURL string) (*ShortenerLinkModel, error) {
	var shortenerLink ShortenerLinkModel
	err := r.db.Where("shortener_url = ?", shortenerURL).First(&shortenerLink).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &shortenerLink, nil
}
