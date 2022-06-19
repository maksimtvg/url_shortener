package services

import (
	"context"
	"errors"
	"log"

	"url_shortener/internal/generator"
	"url_shortener/internal/pkg/shortener"
	"url_shortener/internal/repositories"
)

type UrlShortener struct {
	repo repositories.Repository
	gen  generator.Generator
}

// NewUrlShortener inits paddingIndex for generator.UriGenerator and returns *UrlShortener
// initIndex is equal to the very last primary key in storage.
func NewUrlShortener(r repositories.Repository) *UrlShortener {
	initIndex, err := r.Init()
	if err != nil {
		log.Fatalf("%s", err)
	}

	return &UrlShortener{
		repo: r,
		gen:  generator.NewUriGenerator(initIndex),
	}
}

// Create handles creating new unique URI and inserts it in a storage.
func (u *UrlShortener) Create(ctx context.Context, url *shortener.CreateUrl) (*shortener.UrlResponse, error) {
	response := &shortener.UrlResponse{
		Short: "",
		Views: 0,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		uniqueUrl, err := u.repo.Insert(ctx, url, u.gen)
		if err != nil {
			return response, errors.New(err.Error())
		}

		response.Short = uniqueUrl
		response.Views = 0
		return response, nil
	}
}

// Delete handles removing uri from a storage.
func (u UrlShortener) Delete(ctx context.Context, url *shortener.DeleteUrl) (*shortener.DeleteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		err := u.repo.Delete(ctx, url.GetUrl())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		response := &shortener.DeleteResponse{
			Status: "ok",
		}
		return response, nil
	}
}

// Get searches url model from the storage.
func (u UrlShortener) Get(ctx context.Context, url *shortener.GetUrl) (*shortener.UrlResponse, error) {
	response := &shortener.UrlResponse{
		Short: "",
		Views: 0,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		uri, err := u.repo.Find(ctx, url.GetUrl())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		response.Short = uri.GetShort()
		response.Views = uri.GetViews()
		return response, nil
	}
}

// Redirect searches url model by shortUri and return.
func (u UrlShortener) Redirect(ctx context.Context, url *shortener.RedirectUrl) (*shortener.RedirectResponse, error) {
	response := &shortener.RedirectResponse{
		Url: "",
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		uri, err := u.repo.Find(ctx, url.GetUrl())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		response.Url = uri.GetLong()
		return response, nil
	}
}
