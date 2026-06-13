package domain

type Role string

const (
	RoleOwner    Role = "OWNER"
	RoleManager  Role = "MANAGER"
	RoleEmployee Role = "EMPLOYEE"
)
