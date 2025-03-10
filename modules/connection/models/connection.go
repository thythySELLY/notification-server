package models

import (
	"time"
)

type Connection struct {
	ID                           string    `bson:"_id,omitempty" json:"_id"`
	CreatedAt                    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt                    time.Time `bson:"updatedAt" json:"updatedAt"`
	Status                       string    `bson:"status" json:"status"`
	WebviewServerApiKey          string    `bson:"webviewServerApiKey" json:"webviewServerApiKey"`
	UserDeliveryServerApiKey     string    `bson:"userDeliveryServerApiKey" json:"userDeliveryServerApiKey"`
	WebviewServerId              string    `bson:"webviewServerId" json:"webviewServerId"`
	UserDeliveryServerId         string    `bson:"userDeliveryServerId" json:"userDeliveryServerId"`
	UserDeliveryServerWebHookUrl string    `bson:"userDeliveryServerWebHookUrl" json:"userDeliveryServerWebHookUrl"`
}
