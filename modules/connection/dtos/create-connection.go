package dto

type CreateConnection struct {
	UserDeliveryServerId string `json:"userDeliveryServerId"`
	WebviewServerId      string `json:"webviewServerId"`
	UserDeliveryServerWebHookUrl string `json:"userDeliveryServerWebHookUrl"` 
}
