package shortlink

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"

	"github.com/xcurvnubaim/njajal-gin-golang/internal/modules/common"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/e"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/query"
)

type IUseCase interface {
	CreateShortenerLink(data *CreateShortenerLinkRequestDTO) (*CreateShortenerLinkResponseDTO, e.ApiError)
	GetOriginalURL(shortenerURL string) (*string, e.ApiError)
	GetAllShortenerLink(queryParam *query.QueryParams) (*common.PaginationResponseDTO[GetAllShortenerLinksResponseDTO], e.ApiError)
}

type useCase struct {
	repository IRepository
}

func NewuseCase(repository IRepository) *useCase {
	return &useCase{repository}
}

func (uc *useCase) CreateShortenerLink(data *CreateShortenerLinkRequestDTO) (*CreateShortenerLinkResponseDTO, e.ApiError) {
	if data.ShortenerURL == "" {
		var maxRetry = 5
		for i := 0; i < maxRetry; i++ {
			data.ShortenerURL = uc.GenerateRandomShortenerURL(6)
			check, _ := uc.repository.GetShortenerLinkByShortenerURL(data.ShortenerURL)
			if check == nil {
				break
			}
			if i == maxRetry-1 {
				log.Println("Max retry generate shorten url reached")
				return nil, e.NewApiError(500, "Failed to generate shortener URL")
			}
		}
	} else {
		check, _ := uc.repository.GetShortenerLinkByShortenerURL(data.ShortenerURL)
		if check != nil {
			return nil, e.NewApiError(400, "Shortener URL already exists")
		}
	}

	shortenerLinkModel := NewShortenerLink(data.OriginalURL, data.ShortenerURL)
	err := uc.repository.CreateShortenerLink(shortenerLinkModel)
	if err != nil {
		return nil, e.NewApiError(500, err.Error())
	}

	return &CreateShortenerLinkResponseDTO{
		OriginalURL:  shortenerLinkModel.OriginalURL,
		ShortenerURL: shortenerLinkModel.ShortenerURL,
	}, nil
}

func (uc *useCase) GenerateRandomShortenerURL(length int) string {
	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	// Encode bytes to a base64 string and clean it for URL use
	shortURL := base64.URLEncoding.EncodeToString(bytes)
	shortURL = strings.TrimRight(shortURL, "=") // Remove padding
	shortURL = shortURL[:length]                // Trim to the desired length

	return shortURL
}

func (uc *useCase) GetOriginalURL(shortenerURL string) (*string, e.ApiError) {
	shortenerLink, err := uc.repository.GetShortenerLinkByShortenerURL(shortenerURL)
	if err != nil {
		return nil, e.NewApiError(400, "Shortener URL not found")
	}

	return &shortenerLink.OriginalURL, nil
}

func (uc *useCase) GetAllShortenerLink(queryParam *query.QueryParams) (*common.PaginationResponseDTO[GetAllShortenerLinksResponseDTO], e.ApiError) {
	shortenerLinks, err := uc.repository.GetAllShortenerLink(queryParam.ApplyQuery)
	if err != nil {
		return nil, e.NewApiError(500, err.Error())
	}

	var response common.PaginationResponseDTO[GetAllShortenerLinksResponseDTO]
	data := make([]GetShortenerLink, 0)
	for _, shortenerLink := range shortenerLinks {
		data = append(data, GetShortenerLink{
			ID: 		 shortenerLink.ID.String(),
			OriginalURL:  shortenerLink.OriginalURL,
			ShortenerURL: shortenerLink.ShortenerURL,
			CreatedAt:    shortenerLink.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalCount, err := uc.repository.CountShortenerLink(queryParam.ApplyQuery)
	if err != nil {
		return nil, e.NewApiError(500, err.Error())
	}

	response.Data = &GetAllShortenerLinksResponseDTO{
		ShortenerLink: data,
	}

	response.Meta = queryParam.NewPaginationMeta(int(totalCount))

	return &response, nil
}