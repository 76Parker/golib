package pglib

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mast-se/go-lib/health"
)

type HealthChecker struct {
	pool *pgxpool.Pool
}

func NewHealthChecker(db *pgxpool.Pool) health.Checker {
	return &HealthChecker{
		pool: db,
	}
}

func (r *HealthChecker) Check(ctx context.Context) health.Status {
	if err := r.pool.Ping(ctx); err != nil {
		return health.Status{
			System:   health.PostgresSystemType,
			IsReady:  false,
			ErrorMsg: err.Error(),
		}
	}
	return health.Status{
		System:   health.PostgresSystemType,
		IsReady:  true,
		ErrorMsg: "",
	}
}
