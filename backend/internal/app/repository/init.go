package repository

import "pos-system/backend/internal/app/port"

func New(sql port.Repositories) *port.Repository {
	return &port.Repository{SQL: sql}
}
