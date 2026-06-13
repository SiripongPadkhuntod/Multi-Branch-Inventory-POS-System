package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type APIError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	HTTPStatus  int    `json:"-"`
}

func BindReqJson200Resp[Request any, Response any](
	c *gin.Context,
	f func(ctx context.Context, request Request) (Response, error),
) {
	var req Request
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{Code: "INVALID_REQUEST", Description: err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{Code: "INVALID_REQUEST", Description: err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil && c.Request.Method != http.MethodGet && c.Request.Method != http.MethodDelete {
		c.JSON(http.StatusBadRequest, APIError{Code: "INVALID_REQUEST", Description: err.Error()})
		return
	}
	resp, err := f(injectHeaders(c), req)
	if err != nil {
		c.JSON(statusFromError(err), APIError{Code: "ERROR", Description: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func statusFromError(err error) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return http.StatusRequestTimeout
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusNotFound
	}
	switch err.Error() {
	case "invalid email or password", "invalid refresh token", "missing refresh token", "user not found",
		"missing authorization header", "invalid authorization header", "invalid access token", "user is inactive",
		"authenticated user is required":
		return http.StatusUnauthorized
	case "email and password are required", "barcode is required", "invalid uuid",
		"branch_id is required", "product_id is required",
		"quantity_delta is required", "receive quantity must be greater than zero",
		"stock cannot become negative":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func injectHeaders(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			ctx = context.WithValue(ctx, k, v[0])
		}
	}
	return ctx
}
