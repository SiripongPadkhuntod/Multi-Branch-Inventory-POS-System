package dto

type EmptyStruct struct{}

type SuccessResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
