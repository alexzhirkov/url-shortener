package get

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

type httpHandler struct {
	repo UrlGetter
}

func (h httpHandler) New(ctx *gin.Context, repo UrlGetter) {
	response, err := h.repo.GetUrl(ctx.Param("alias"))
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s", err)})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"id": response})
}
