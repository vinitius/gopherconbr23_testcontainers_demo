package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type AccountModel struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            string `bun:"account_id,pk" json:"account_id"`
	Type          string `bun:"type" json:"type"`
	Label         string `bun:"label" json:"label"`
	Currency      string `bun:"currency" json:"currency"`
}

type AccountsSQLRepository struct {
	db bun.IDB
}

func NewAccountsSQLRepository(db bun.IDB) AccountsSQLRepository {
	return AccountsSQLRepository{db: db}
}

func (r AccountsSQLRepository) Create(ctx context.Context, account AccountModel) error {
	_, err := r.db.NewInsert().Model(&account).Exec(ctx)

	return err
}

func (r AccountsSQLRepository) FindAll(ctx context.Context) ([]AccountModel, error) {
	var accounts []AccountModel
	query := r.db.NewSelect().Model(&accounts)
	err := query.Scan(ctx)

	return accounts, err
}
