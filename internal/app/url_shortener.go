package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"url_shortener/internal/pkg/shortener"
)

type UrlShortener struct {
	db *pgxpool.Pool
}

func NewUrlShortener(db *pgxpool.Pool) *UrlShortener {
	return &UrlShortener{
		db: db,
	}
}

func (u UrlShortener) Create(ctx context.Context, url *shortener.CreateUrl) (*shortener.UrlResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UrlShortener) Delete(ctx context.Context, url *shortener.DeleteUrl) (*shortener.DeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UrlShortener) Get(ctx context.Context, url *shortener.GetUrl) (*shortener.UrlResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UrlShortener) Redirect(ctx context.Context, url *shortener.RedirectUrl) (*shortener.RedirectResponse, error) {
	//TODO implement me
	panic("implement me")
}
