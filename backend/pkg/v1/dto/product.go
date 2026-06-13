package dto

import "pos-system/backend/pkg/v1/model"

type ProductListRequest struct {
	Query  string `form:"q"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type ProductBarcodeRequest struct {
	Barcode string `uri:"barcode"`
}

type ProductListResponse struct {
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Data        []model.Product `json:"data"`
}

type ProductResponse struct {
	Code        string        `json:"code"`
	Description string        `json:"description"`
	Data        model.Product `json:"data"`
}
