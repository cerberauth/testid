package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	challenge, exists := c.GetQuery("logout_challenge")
	if !exists || challenge == "" {
		c.String(http.StatusBadRequest, "Logout challenge not found")
		return
	}

	acceptResp, r, err := h.hydraApi.OAuth2API.AcceptOAuth2LogoutRequest(c).LogoutChallenge(challenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLogoutRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
}
