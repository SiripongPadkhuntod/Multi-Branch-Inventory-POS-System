package handler

import (
	"context"
	"errors"
	"strings"

	"pos-system/backend/internal/app/domain"
	"pos-system/backend/pkg/v1/dto"
	"pos-system/backend/pkg/v1/model"
)

func (h *handler) ProductListHandler(ctx context.Context, req dto.ProductListRequest) (*dto.ProductListResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	products, err := h.svc.ListProducts(ctx, strings.TrimSpace(req.Query), req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	data := make([]model.Product, 0, len(products))
	for _, product := range products {
		data = append(data, productModel(product))
	}
	return &dto.ProductListResponse{Code: "SUCCESS", Description: "success", Data: data}, nil
}

func (h *handler) ProductBarcodeHandler(ctx context.Context, req dto.ProductBarcodeRequest) (*dto.ProductResponse, error) {
	if strings.TrimSpace(req.Barcode) == "" {
		return nil, errors.New("barcode is required")
	}
	product, err := h.svc.ProductByBarcode(ctx, strings.TrimSpace(req.Barcode))
	if err != nil {
		return nil, err
	}
	return &dto.ProductResponse{Code: "SUCCESS", Description: "success", Data: productModel(*product)}, nil
}

func productModel(product domain.Product) model.Product {
	return model.Product{
		ID:          product.ID.String(),
		SKU:         product.SKU,
		Barcode:     product.Barcode,
		QRCode:      product.QRCode,
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  product.CategoryID.String(),
		ImageURL:    product.ImageURL,
		CostPrice:   product.CostPrice,
		SellPrice:   product.SellPrice,
		Status:      product.Status,
	}
}
