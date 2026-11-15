package models

type CreateRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	Photos      []string
}

type Advertisement struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Photos      []Photo `json:"photos"`
}

type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type AdvertisementList struct {
	Items []Advertisement `json:"items"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
