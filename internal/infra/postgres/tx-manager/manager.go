package txmanager

import (
	"context"

	"example/internal/domain"
	"example/internal/infra/postgres"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	pool postgres.PgxPool
}

func NewRepository(pool postgres.PgxPool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.DoWithOptions(ctx, domain.TxOptions{
		IsolationLevel: domain.IsolationDefault,
		ReadOnly:       false,
	}, fn)
}

func (r *Repository) DoWithOptions(ctx context.Context, opts domain.TxOptions, fn func(ctx context.Context) error) error {
	isoLevel := convertIsolationLevel(opts.IsolationLevel)
	txOpts := pgx.TxOptions{
		IsoLevel:   isoLevel,
		AccessMode: pgx.ReadWrite,
	}
	if opts.ReadOnly {
		txOpts.AccessMode = pgx.ReadOnly
	}
	tx, err := r.pool.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func convertIsolationLevel(level domain.IsolationLevel) pgx.TxIsoLevel {
	switch level {
	case domain.IsolationReadUncommitted:
		return pgx.ReadUncommitted
	case domain.IsolationReadCommitted:
		return pgx.ReadCommitted
	case domain.IsolationRepeatableRead:
		return pgx.RepeatableRead
	case domain.IsolationSerializable:
		return pgx.Serializable
	default:
		return pgx.ReadCommitted
	}
}
