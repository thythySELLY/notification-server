package domain

import "notification-server/modules/webview-server/models"

type GetWebViewList struct {
	List          []models.WebViewServer `json:"list"`
	NextPageToken string                 `json:"nextPageToken"`
}
