package domain

type ConnectionResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}
