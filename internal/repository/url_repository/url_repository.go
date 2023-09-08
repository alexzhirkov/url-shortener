package url_repository

import "errors"

type ShortUrl struct {
	Id    int
	Url   string
	Alias string
}

var (
	ErrAliasNotFound = errors.New("alias not found")
	ErrUrlNotFound   = errors.New("http not found")
	ErrUrlExists     = errors.New("http exists")
)
