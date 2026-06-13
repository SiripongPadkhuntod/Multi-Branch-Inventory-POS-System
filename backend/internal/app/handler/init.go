package handler

import (
	"context"
	"errors"
	"strings"

	"pos-system/backend/internal/app/port"
	"pos-system/backend/pkg/v1/dto"
)

type handler struct {
	svc port.Service
}

func New(svc port.Service) port.Handler {
	return &handler{svc: svc}
}

func (h *handler) HealthHandler(ctx context.Context, _ dto.EmptyStruct) (*dto.HealthResponse, error) {
	return &dto.HealthResponse{Code: "SUCCESS", Description: "success", Status: "ok"}, nil
}

func (h *handler) LoginHandler(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if strings.TrimSpace(req.Email) == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}
	access, refresh, user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Code: "SUCCESS", Description: "success", AccessToken: access, RefreshToken: refresh, User: user}, nil
}

func (h *handler) RefreshHandler(ctx context.Context, req dto.RefreshRequest) (*dto.LoginResponse, error) {
	access, refresh, user, err := h.svc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Code: "SUCCESS", Description: "success", AccessToken: access, RefreshToken: refresh, User: user}, nil
}

func (h *handler) LogoutHandler(ctx context.Context, req dto.LogoutRequest) (*dto.SuccessResponse, error) {
	if err := h.svc.Logout(ctx, req.RefreshToken); err != nil {
		return nil, err
	}
	return &dto.SuccessResponse{Code: "SUCCESS", Description: "success"}, nil
}
