package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Error(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", gin.H{
		"Error": "An error occurred. Please try again.",
	})
}
