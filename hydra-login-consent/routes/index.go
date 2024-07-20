package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Index(c *gin.Context) {
	c.Header("Cache-Control", "max-age=3600, s-maxage=86400, stale-while-revalidate=86400, stale-if-error=86400")
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
