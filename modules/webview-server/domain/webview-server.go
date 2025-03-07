package domain

type WebViewResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    any `json:"data"`
}
