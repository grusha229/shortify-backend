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

type GetURLStatDataPayload struct {
    LinkId uuid.UUID `json:"link_id"`
}

type GetURLStatDataDB struct {
    ID uuid.UUID `json:"id"`
    LinkId uuid.UUID `json:"link_id"`
    VisitedAt string `json:"visited_at"`
    IPAddress string `json:"ip_address"`
    UserAgent string `json:"user_agent"`
    UtmSource string `json:"utm_source"`
}

type GetURLStatDataResponse struct {
    ID uuid.UUID `json:"id"`
    LinkId uuid.UUID `json:"link_id"`
    VisitedAt string `json:"visited_at"`
    Location string `json:"location"`
    UtmSource string `json:"utm_source"`
}