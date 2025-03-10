package dto

type UpdateUserDelivery struct {
	ID                           string `json:"id"`
	UserDeliveryServerWebHookUrl string `bson:"userDeliveryServerWebHookUrl" json:"userDeliveryServerWebHookUrl"`
}
