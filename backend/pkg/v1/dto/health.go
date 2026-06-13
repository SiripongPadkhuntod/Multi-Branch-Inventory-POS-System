package dto

type HealthResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
