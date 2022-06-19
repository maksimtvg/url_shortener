package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url_shortener/internal/generator"
	"url_shortener/internal/pkg/shortener"
	mock_repositories "url_shortener/internal/repositories/mock"
)

// tests Create method.
func TestUrlShortener_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repositories.NewMockRepository(ctrl)

	// build stubs
	urlRequest := &shortener.CreateUrl{
		Url: "https://example.com/uri",
	}

	g := generator.NewUriGenerator(0)
	uri, _ := g.GenerateUri()

	repo.EXPECT().
		Insert(gomock.Any(), urlRequest, gomock.Any()).
		Times(1).
		Return(uri, nil)

	initialIndex := int64(0)
	service := &UrlShortener{
		repo: repo,
		gen:  generator.NewUriGenerator(initialIndex),
	}

	create, err := service.Create(context.Background(), &shortener.CreateUrl{
		Url: "https://example.com/uri",
	})
	require.NoError(t, err)
	assert.Equal(t, create.Short, uri)
}

// tests Get method.
func TestUrlShortener_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repositories.NewMockRepository(ctrl)

	// build stubs
	response := &shortener.UrlResponse{
		Long:  "https://example.com/uri",
		Short: "6LAze",
		Views: 10,
	}

	repo.EXPECT().
		Find(gomock.Any(), response.GetShort()).
		Times(1).
		Return(response, nil)

	initialIndex := int64(0)
	service := &UrlShortener{
		repo: repo,
		gen:  generator.NewUriGenerator(initialIndex),
	}

	got, err := service.Get(context.Background(), &shortener.GetUrl{Url: "6LAze"})
	require.NoError(t, err)
	assert.Equal(t, got.GetViews(), response.GetViews())
	assert.Equal(t, got.GetShort(), response.GetShort())
}

// tests Redirect method.
func TestUrlShortener_Redirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repositories.NewMockRepository(ctrl)

	url := &shortener.RedirectUrl{
		Url: "6LAze",
	}
	response := &shortener.UrlResponse{
		Long:  "https://example.com/uri",
		Short: "6LAze",
		Views: 10,
	}

	repo.EXPECT().
		Find(gomock.Any(), url.GetUrl()).
		Times(1).
		Return(response, nil)

	initialIndex := int64(0)
	service := &UrlShortener{
		repo: repo,
		gen:  generator.NewUriGenerator(initialIndex),
	}

	got, err := service.Redirect(context.Background(), url)
	require.NoError(t, err)
	assert.Equal(t, got.GetUrl(), response.GetLong())
}
