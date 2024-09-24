package common

type Response[D any] struct {
	Status  bool    `json:"status"`
	Message string  `json:"message"`
	Data    *D      `json:"data"`
	Error   *string `json:"error"`
}
