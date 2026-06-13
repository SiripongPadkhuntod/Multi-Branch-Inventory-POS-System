package domain

import appdomain "pos-system/backend/internal/app/domain"

type Role = appdomain.Role

const (
	RoleOwner    = appdomain.RoleOwner
	RoleManager  = appdomain.RoleManager
	RoleEmployee = appdomain.RoleEmployee
)
