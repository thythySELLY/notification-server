package dto

type GetConnections struct {
	UserDeliveryServerId string `json:"userDeliveryServerId"`
	WebviewServerId      string `json:"webviewServerId"`
	Status               string `form:"status"`
	Limit                int    `form:"limit"`
	PageToken            string `form:"pageToken"`
}
