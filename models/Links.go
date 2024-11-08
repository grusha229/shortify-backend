package LinksModels

import "github.com/google/uuid"

type ShortenRequest struct {
    URL string `json:"url"`
}

type ShortenResponse struct {
    ID uuid.UUID `json:"id"`
    ShortURL string `json:"short_url"`
    OriginalURL string `json:"original_url"`
}