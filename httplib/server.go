package httplib

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Address           string        `yaml:"address"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes"`
}

func NewHTTPServer(ctx context.Context, cfg Config, handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:              cfg.Address,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	return s
}
