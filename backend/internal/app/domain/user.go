package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID   `json:"id"`
	BranchID     *uuid.UUID  `json:"branch_id"`
	BranchIDs    []uuid.UUID `json:"branch_ids,omitempty"`
	Role         Role        `json:"role"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

type EmployeeSalesSummary struct {
	UserID     uuid.UUID `json:"user_id"`
	BranchID   uuid.UUID `json:"branch_id"`
	BranchCode string    `json:"branch_code"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       Role      `json:"role"`
	Status     string    `json:"status"`
	SalesCount int64     `json:"sales_count"`
	Revenue    int64     `json:"revenue"`
}
