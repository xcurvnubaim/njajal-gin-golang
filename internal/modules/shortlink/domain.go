package shortlink

import (
	"github.com/xcurvnubaim/njajal-gin-golang/internal/modules/common"
)

type ShortenerLinkModel struct {
	common.BaseModels
	OriginalURL  string `gorm:"column:original_url;not null"`
	ShortenerURL string `gorm:"column:shortener_url;not null"`
}

func (ShortenerLinkModel) TableName() string {
	return "shortener_links"
}

func NewShortenerLink(originalURL, shortenerURL string) *ShortenerLinkModel {
	return &ShortenerLinkModel{
		BaseModels:   common.NewBaseModels(),
		OriginalURL:  originalURL,
		ShortenerURL: shortenerURL,
	}
}
