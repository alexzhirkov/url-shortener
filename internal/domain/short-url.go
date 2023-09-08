package domain

import (
	"fmt"
	"strings"
)

/*
еще вариант/ странный какой-то. в домене будет dto репозитория. а доступ к полям? публичный? да ну
type ShortUrl struct {
    urlRepo port.UrlRepository
}

func NewCar(urlRepo port.UrlRepository) &ShortUrl {
    return &ShortUrl{
        urlRepo: urlRepo,
    }
}
*/

type ShortUrl struct {
	id    int
	url   string
	alias string
}

func NewShortUrl(id int, url string, alias string) (*ShortUrl, error) {
	if id < 1 {
		return nil, fmt.Errorf("%w: shorturl id must be greater than zero", ErrRequired)
	}
	if url == "" {
		return nil, fmt.Errorf("%w: url is required", ErrRequired)
	}
	if strings.Index(url, "http") != 0 {
		return nil, fmt.Errorf("%w: url must starts with 'http'", ErrRequired)
	}

	if alias == "" {
		return nil, fmt.Errorf("%w: shorturl alias is required", ErrRequired)
	}
	return &ShortUrl{
		id:    id,
		url:   url,
		alias: alias,
	}, nil
}

func (u *ShortUrl) GetId() int {
	return u.id
}
func (u *ShortUrl) GetUrl() string {
	return u.url
}
func (u *ShortUrl) GetAlias() string {
	return u.alias
}
