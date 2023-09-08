package in_memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexzhirkov/url-shortener/internal/domain"
	"github.com/alexzhirkov/url-shortener/internal/repository/url_repository"
	"sync"
)

type Storage struct {
	data map[string]*url_repository.ShortUrl
	//data map[int]*ShortUrl
	mu sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[string]*url_repository.ShortUrl),
	}
}

func (store *Storage) GetUrl(ctx context.Context, alias string) (*domain.ShortUrl, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if value, ok := store.data[alias]; ok {
		storeUrl, err := domain.NewShortUrl(value.Id, value.Url, value.Alias)
		//err := json.Unmarshal(value, &url)
		if err != nil {
			return nil, errors.New("can't create short url")
		}
		//todo:make copy of store value/ smth like url, err := urlStoreToDomain(storeUrl)
		return storeUrl, nil
	}

	return nil, url_repository.ErrAliasNotFound
}

// Count returns the number of ports
func (store *Storage) Count() (int, error) {
	return len(store.data), nil
}

// CreateUrl creates or updates a port
func (store *Storage) CreateUrl(ctx context.Context, shortUrl *domain.ShortUrl) error {

	if _, ok := store.data[shortUrl.GetAlias()]; ok {
		return fmt.Errorf("alias already exists: %v", shortUrl.GetAlias())
	}
	store.mu.Lock()
	store.data[shortUrl.GetAlias()] = &url_repository.ShortUrl{
		Id:    shortUrl.GetId(),
		Url:   shortUrl.GetUrl(),
		Alias: shortUrl.GetAlias(),
	}
	store.mu.Unlock()
	return nil
}
