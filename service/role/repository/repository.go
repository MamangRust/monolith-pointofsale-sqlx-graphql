package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

func NewRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		RoleCommand: NewRoleCommandRepository(DB),
		RoleQuery:   NewRoleQueryRepository(DB),
	}
}
