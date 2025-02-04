package shortlink

type (
	CreateShortenerLinkRequestDTO struct {
		OriginalURL  string `json:"original_url" binding:"url,required"`
		ShortenerURL string `json:"shortener_url" `
	}

	CreateShortenerLinkResponseDTO struct {
		OriginalURL  string `json:"original_url"`
		ShortenerURL string `json:"shortener_url"`
	}

	GetShortenerLink struct {
		ID           string `json:"id"`
		OriginalURL  string `json:"original_url"`
		ShortenerURL string `json:"shortener_url"`
		CreatedAt    string `json:"created_at"`
	}

	GetAllShortenerLinksResponseDTO struct {
		ShortenerLink []GetShortenerLink `json:"shortener_links"`
	}
)
