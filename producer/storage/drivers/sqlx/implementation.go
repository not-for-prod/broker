package sqlx

import (
	"context"

	"github.com/Masterminds/squirrel"
	trm "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/not-for-prod/broker/producer"
)

var _ producer.Storage = &Implementation{}
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Implementation struct {
	db        *sqlx.DB
	ctxGetter *trm.CtxGetter
}

func NewImplementation(db *sqlx.DB, ctxGetter *trm.CtxGetter) *Implementation {
	return &Implementation{
		db:        db,
		ctxGetter: ctxGetter,
	}
}

func (i *Implementation) tr(ctx context.Context) trm.Tr {
	return i.ctxGetter.DefaultTrOrDB(ctx, i.db)
}
