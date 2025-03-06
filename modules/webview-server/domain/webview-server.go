package domain

import "notification-server/modules/webview-server/models"

// WebViewResponse chứa dữ liệu trả về client
type WebViewResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    WebViewList `json:"data"`
}

type WebViewList struct {
	List          []models.WebViewServer `json:"list"`
	NextPageToken string           `json:"nextPageToken"`
}
