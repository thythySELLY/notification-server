package dto

type GetUserDeliveryList struct {
	Keyword   string `form:"keyword"`
	Status    string `form:"status"`
	Limit     int    `form:"limit"`
	PageToken string `form:"pageToken"`
}
