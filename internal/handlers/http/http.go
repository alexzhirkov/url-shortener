package http

import (
	"context"
	"github.com/alexzhirkov/url-shortener/internal/domain"
	"github.com/gin-gonic/gin"
)

type ShortUrl struct {
	Id    int
	Url   string
	Alias string
}

type urlUseCase interface {
	GetUrl(ctx context.Context, alias string) (*domain.ShortUrl, error)
	Count() (int, error)
	CreateUrl(ctx context.Context, shorturl *domain.ShortUrl) error
}

type Handler struct {
	urlUseCase urlUseCase
}

func NewHTTPHandler(urlUseCase urlUseCase) *Handler {
	return &Handler{
		urlUseCase: urlUseCase,
	}
}

func (handler *Handler) GetUrl(ctx *gin.Context) {
	alias := ctx.Param("alias")
	url, err := handler.urlUseCase.GetUrl(ctx, alias)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	//fmt.Println(*url)
	ctx.JSON(200, gin.H{
		"id":    url.GetId(),
		"alias": url.GetAlias(),
		"url":   url.GetUrl(),
	})
}

func (handler *Handler) CreateUrl(ctx *gin.Context) {
	var newUrl ShortUrl
	if err := ctx.BindJSON(&newUrl); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
	}
	//fmt.Println(newUrl)
	shortUrl, err := domain.NewShortUrl(newUrl.Id, newUrl.Url, newUrl.Alias)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}
	if err := handler.urlUseCase.CreateUrl(ctx, shortUrl); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"success": true})
}

func (handler *Handler) CountUrls(ctx *gin.Context) {
	count, err := handler.urlUseCase.Count()
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"count": count})
}
