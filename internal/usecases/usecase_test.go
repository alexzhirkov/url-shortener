package usecases_test

import (
	"github.com/alexzhirkov/url-shortener/internal/usecases"
	usecases2 "github.com/alexzhirkov/url-shortener/internal/usecases"
	"github.com/alexzhirkov/url-shortener/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_CountUrls(t *testing.T) {
	type testCase struct {
		name      string
		count     int
		mockError error
	}

	test_cases := []testCase{
		{
			name:  "error",
			count: 3,
		},
		{
			name:  "count equal 0",
			count: 0,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			urlRepositoryMock := mocks.NewUrlRepository(t)

			//if tc.respError == "" || tc.mockError != nil {
			urlRepositoryMock.
				On("Count").
				Return(tc.count, tc.mockError).
				Once()
			/*urlRepositoryMock.
			On("CreateUrl", context.Background(), &domain.ShortUrl{Id: 1, Url: "test", Alias: "test12"}).
			Return(nil).
			Once()*/
			//}

			service, _ := usecases2.New(usecases.WithUrlRepository(urlRepositoryMock))
			//_ = service.CreateUrl(context.Background(), &domain.ShortUrl{Id: 1, Url: "test", Alias: "test12"})

			count, _ := service.Count()
			assert.Equal(t, tc.count, count)

		})
	}
}
