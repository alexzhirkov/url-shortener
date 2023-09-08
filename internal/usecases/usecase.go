package usecases

import (
	"context"
	"fmt"
	"github.com/alexzhirkov/url-shortener/internal/domain"
)

// UrlRepository is an HTTP repository for the http service
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=urlRepository
type urlRepository interface {
	CreateUrl(ctx context.Context, shorturl *domain.ShortUrl) error
	Count() (int, error)
	GetUrl(ctx context.Context, alias string) (*domain.ShortUrl, error)
}

type UrlUseCaseConfiguration func(uuc *UrlUseCase) error

// UrlUseCase is an HTTP service
type UrlUseCase struct {
	repo urlRepository
}

// New creates a new http service
func New(cfgs ...UrlUseCaseConfiguration) (*UrlUseCase, error) {
	uuc := &UrlUseCase{}
	for _, cfg := range cfgs {
		err := cfg(uuc)
		if err != nil {
			return nil, err
		}
	}
	return uuc, nil
}

func WithUrlRepository(ur urlRepository) UrlUseCaseConfiguration {
	return func(uuc *UrlUseCase) error {
		uuc.repo = ur
		return nil
	}
}

/*func WithMemoryRepository() UrlUseCaseConfiguration {
	urlStore := in_memory.New()
	return WithUrlRepository(urlStore)
}*/

// Get returns a port by id
func (s UrlUseCase) GetUrl(ctx context.Context, alias string) (*domain.ShortUrl, error) {
	return s.repo.GetUrl(ctx, alias)
}

// CountUrls returns the number of ports
func (s UrlUseCase) Count() (int, error) {
	fn := "usecase.Count"
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("%s, db.Prepare: %w", fn, err)
	}
	return count, nil
}

// CreateUrl creates or updates a port
func (s UrlUseCase) CreateUrl(ctx context.Context, shorturl *domain.ShortUrl) error {
	return s.repo.CreateUrl(ctx, shorturl)
}
