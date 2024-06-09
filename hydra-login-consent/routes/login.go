package routes

import (
	"fmt"
	"net/http"
	"os"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go"
)

func newAcceptLoginRequest(subject string) *hydraClient.AcceptLoginRequest {
	acceptLoginRequest := hydraClient.NewAcceptLoginRequest(subject)
	acceptLoginRequest.SetRemember(true)
	acceptLoginRequest.SetRememberFor(3600 * 12)
	return acceptLoginRequest
}

func (h *Handler) Login(c *gin.Context) {
	challenge, exists := c.GetQuery("login_challenge")
	if !exists || challenge == "" {
		c.String(http.StatusBadRequest, "Login challenge not found")
		return
	}

	loginRequest, r, err := h.hydraApi.AdminApi.GetLoginRequest(c).LoginChallenge(challenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetLoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	if loginRequest.Skip {
		fmt.Printf("Accepting login request because it was skipped\n")

		acceptLoginRequest := newAcceptLoginRequest(loginRequest.GetSubject())
		acceptResp, r, err := h.hydraApi.AdminApi.AcceptLoginRequest(c).LoginChallenge(challenge).AcceptLoginRequest(*acceptLoginRequest).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLoginRequest``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
		return
	}

	loginHint := "john.doe@example.com"
	if loginRequest.GetOidcContext().LoginHint != nil {
		loginHint = *loginRequest.GetOidcContext().LoginHint
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Challenge": challenge,
		"LoginHint": loginHint,
	})
}

type PostLoginForm struct {
	Challenge string `form:"challenge" binding:"required"`
	Email     string `form:"email" binding:"required"`
}

func (h *Handler) PostLogin(c *gin.Context) {
	var form PostLoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	hash := sha256.Sum256([]byte(form.Email))
	subject := hex.EncodeToString(hash[:])

	acceptLoginRequest := newAcceptLoginRequest(subject)
	acceptLoginRequest.SetContext(map[string]interface{}{
		"email": form.Email,
	})
	acceptResp, r, err := h.hydraApi.AdminApi.AcceptLoginRequest(c).LoginChallenge(form.Challenge).AcceptLoginRequest(*acceptLoginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
}
