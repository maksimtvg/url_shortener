// Package repositories
package repositories

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"url_shortener/internal/generator"
	"url_shortener/internal/pkg/shortener"
)

type DBRepository struct {
	Pg *pgxpool.Pool
}

func NewDBRepository(dbPool *pgxpool.Pool) *DBRepository {
	return &DBRepository{Pg: dbPool}
}

// Insert inserts new urls
func (r *DBRepository) Insert(ctx context.Context, url *shortener.CreateUrl, gen generator.Generator) (string, error) {
	var (
		err       error
		tx        pgx.Tx
		uniqueUri string
	)

	tx, err = r.Pg.Begin(ctx)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	uniqueUri, err = gen.GenerateUri()
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}
	_, err = tx.Exec(ctx, "INSERT INTO urls (url, short_uri, views) VALUES($1, $2, $3);",
		url.GetUrl(),
		uniqueUri,
		0,
	)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return "", status.Error(codes.Internal, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}
	return uniqueUri, nil
}

// Delete add deleted_att for cortege
func (r *DBRepository) Delete(ctx context.Context, url string) error {
	var (
		tx pgx.Tx
	)

	tx, err := r.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE urls SET deleted_at=$1 WHERE short_uri=$2", time.Now(), url)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	err = tx.Commit(ctx)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

// Init gets max id in database for initiating generator
func (r *DBRepository) Init() (int64, error) {
	var (
		ctx context.Context
		num int
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := r.Pg.QueryRow(ctx, "SELECT COALESCE(max(id),0) from urls;").Scan(&num)
	if err != nil {
		return -1, status.Error(codes.Internal, err.Error())
	}
	return int64(num), nil
}

// Find searches by shortUri string
func (r *DBRepository) Find(ctx context.Context, shortUri string) (*shortener.UrlResponse, error) {
	var (
		long  pgtype.Text
		short pgtype.Text
		views int64
	)

	err := r.Pg.QueryRow(ctx, "SELECT url, short_uri, views FROM urls WHERE short_uri=$1", shortUri).
		Scan(&long, &short, &views)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	urlRsp := &shortener.UrlResponse{
		Long:  long.String,
		Short: short.String,
		Views: views,
	}
	return urlRsp, nil
}
