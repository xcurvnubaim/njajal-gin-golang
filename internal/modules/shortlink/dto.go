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
)
