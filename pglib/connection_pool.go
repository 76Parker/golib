package pglib

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNilConfig = errors.New("postgres: nil config")
)

type Config struct {
	Host       string     `yaml:"host" validate:"required"`
	Port       int        `yaml:"port" validate:"required"`
	Database   string     `yaml:"database" validate:"required"`
	SSL        string     `yaml:"ssl" validate:"required"`
	Username   string     `yaml:"username" validate:"required"`
	Password   string     `yaml:"password" validate:"required"`
	PoolConfig PoolConfig `yaml:"pool_config" validate:"required"`
}

type PoolConfig struct {
	MaxConnections        int           `yaml:"max_connections" validate:"required"`
	MinConnections        int           `yaml:"min_connections" validate:"required"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime" validate:"required"`
	MaxConnIdleTime       time.Duration `yaml:"max_conn_idle_time" validate:"required"`
	HealthCheckPeriod     time.Duration `yaml:"health_check_period" validate:"required"`
	ConnectTimeout        time.Duration `yaml:"connect_timeout" validate:"required"`
}

func (c *Config) DSN() string {
	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.Username, c.Password, c.Database, c.SSL)
	return connURL
}

func New(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}
	c, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	c.MaxConns = int32(cfg.PoolConfig.MaxConnections)
	c.MinConns = int32(cfg.PoolConfig.MinConnections)
	c.MaxConnLifetime = cfg.PoolConfig.MaxConnectionLifetime
	c.MaxConnIdleTime = cfg.PoolConfig.MaxConnIdleTime
	c.HealthCheckPeriod = cfg.PoolConfig.HealthCheckPeriod
	c.ConnConfig.ConnectTimeout = cfg.PoolConfig.ConnectTimeout
	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
