package domain

import "notification-server/modules/connection/models"

type GetUserDeliveryList struct {
	List          []models.Connection `json:"list"`
	NextPageToken string              `json:"nextPageToken"`
}
