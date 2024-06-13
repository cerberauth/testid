package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Error(c *gin.Context) {
	errorTitle, errorTitleExists := c.GetQuery("error")
	if !errorTitleExists {
		errorTitle = "An error occurred"
	}

	errorDescription, errorDescriptionExists := c.GetQuery("error_description")
	if !errorDescriptionExists {
		errorDescription = ""
	}

	c.HTML(http.StatusOK, "error.html", gin.H{
		"ErrorTitle":       errorTitle,
		"ErrorDescription": errorDescription,
	})
}
