package repository

import (
	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user_role"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        roleadapter.QueryRepository
	UserRole    userroleadapter.CommandRepository
}

type GuardOptions struct {
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

type Deps struct {
	Db *sqlx.DB

	RoleQueryClient pbrole.RoleQueryServiceClient
	UserRoleClient  pbuserrole.UserRoleServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.Db),
		UserQuery:   NewUserQueryRepository(deps.Db),
		Role:        roleadapter.New(deps.RoleQueryClient, g.Role...),
		UserRole:    userroleadapter.New(deps.UserRoleClient, g.UserRole...),
	}
}
