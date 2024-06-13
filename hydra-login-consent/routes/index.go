package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
