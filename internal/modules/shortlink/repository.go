package shortlink

import (
	"log"

	"gorm.io/gorm"
)

type IRepository interface {
	CreateShortenerLink(data *ShortenerLinkModel) error
	GetShortenerLinkByShortenerURL(shortenerURL string) (*ShortenerLinkModel, error)
	CountShortenerLink(applyQuery func(*gorm.DB) *gorm.DB) (int64, error)
	GetAllShortenerLink(applyQuery func(*gorm.DB) *gorm.DB) ([]*ShortenerLinkModel, error)
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

func (r *repository) CountShortenerLink(applyQuery func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	err := applyQuery(r.db).Model(&ShortenerLinkModel{}).Count(&count).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

func (r *repository) GetAllShortenerLink(applyQuery func(*gorm.DB) *gorm.DB) ([]*ShortenerLinkModel, error) {
	var shortenerLinks []*ShortenerLinkModel
	err := applyQuery(r.db).Find(&shortenerLinks).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return shortenerLinks, nil
}
