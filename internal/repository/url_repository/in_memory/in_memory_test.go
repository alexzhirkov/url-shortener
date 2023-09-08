package in_memory

import (
	"context"
	"errors"
	"github.com/alexzhirkov/url-shortener/internal/domain"
	"github.com/alexzhirkov/url-shortener/internal/repository/url_repository"
	"testing"
)

func TestMemoryRepository_getUrl(t *testing.T) {
	type testCase struct {
		name          string
		alias         string
		expectedError error
	}

	surl, err := domain.NewShortUrl(1, "some url", "test_alias")
	if err != nil {
		t.Errorf("error creating domain url %v", err)
	}

	repo := New()
	repo.CreateUrl(context.Background(), surl)

	testCases := []testCase{
		{
			name:          "no url by alias",
			alias:         "test",
			expectedError: url_repository.ErrAliasNotFound,
		},
		{
			name:          "url found",
			alias:         "test_alias",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.GetUrl(context.Background(), tc.alias)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error '%v', got '%v'", tc.expectedError, err)
			}
		})
	}

}
