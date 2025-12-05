package request

type LogRequest struct {
	Data  string `json:"data" validate:"required"`
	Label string `json:"label" validate:"required"`
	Token string `json:"token" validate:"required,uuid"`
}
