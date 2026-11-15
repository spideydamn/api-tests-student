package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

type AddRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type Favorite struct {
	FavoriteID    string                            `json:"favorite_id"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

type FavoriteList struct {
	Items []Favorite `json:"items"`
	Total int        `json:"total"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
