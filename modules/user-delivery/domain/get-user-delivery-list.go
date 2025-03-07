package domain

import "notification-server/modules/user-delivery/models"

type GetUserDeliveryList struct {
	List          []models.UserDelivery `json:"list"`
	NextPageToken string                 `json:"nextPageToken"`
}
