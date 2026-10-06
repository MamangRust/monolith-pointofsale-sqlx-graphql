package repository

import (
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	useradapter "github.com/MamangRust/monolith-pointofsale-pkg/adapter/user"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               useradapter.QueryRepository
}

type GuardOptions struct {
	User []adapter.GuardOption
}

func NewRepositories(
	DB *sqlx.DB,
	userClient pbuser.UserQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB),
		MerchantCommand:         NewMerchantCommandRepository(DB),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB),
		UserQuery:               useradapter.NewQueryAdapter(userClient, g.User...),
	}
}
