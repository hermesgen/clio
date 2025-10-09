package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var key = hm.Key

type ClioRepo struct {
	*hm.BaseRepo
	db *sqlx.DB
}

func NewClioRepo(qm *hm.QueryManager, params hm.XParams) *ClioRepo {
	return &ClioRepo{
		BaseRepo: hm.NewRepo("sqlite-auth-repo", qm, params),
	}
}

// Setup the database connection.
func (repo *ClioRepo) Setup(ctx context.Context) error {
	dsn, ok := repo.Cfg().StrVal(key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
	}
	repo.db = db
	return nil
}

// Stop closes the database connection.
func (repo *ClioRepo) Stop(ctx context.Context) error {
	if repo.db != nil {
		return repo.db.Close()
	}
	return nil
}

// DB returns the underlying *sqlx.DB for transaction management.
func (repo *ClioRepo) DB() *sqlx.DB {
	return repo.db
}

// getExec returns the correct Execer (Tx or DB) from context.
func (repo *ClioRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := hm.TxFromContext(ctx)
	if ok {
		if sqlxTx, ok := tx.(*sqlx.Tx); ok {
			return sqlxTx
		}
	}
	return repo.db
}

func (r *ClioRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, nil, err
	}
	ctxWithTx := hm.WithTx(ctx, tx)
	return ctxWithTx, tx, nil
}
