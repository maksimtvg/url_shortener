package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
)

//newPoolConfig parses config and returns new pgs pool
func newPoolConfig(cfg *DBConfig) (*pgxpool.Config, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.DBUserName),
		url.QueryEscape(cfg.DBPassword),
		cfg.Host,
		cfg.DBPort,
		cfg.DBName,
		cfg.Timeout,
	)

	pool, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// newConnection creates new Postgres DB connection
func newConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Connect return new DB connection
func Connect(dbConfig *DBConfig) (*pgxpool.Pool, error) {
	poolConfig, err := newPoolConfig(dbConfig)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = MaxConn
	c, err := newConnection(poolConfig)
	if err != nil {
		return nil, err
	}
	return c, nil
}
