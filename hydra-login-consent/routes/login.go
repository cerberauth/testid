package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go/v2"
)

func newAcceptLoginRequest(subject string) *hydraClient.AcceptOAuth2LoginRequest {
	acceptLoginRequest := hydraClient.NewAcceptOAuth2LoginRequest(subject)
	acceptLoginRequest.SetRemember(true)
	acceptLoginRequest.SetRememberFor(3600 * 12)
	return acceptLoginRequest
}

func extractNameFromEmail(email string) (string, string, string) {
	// Assuming email format is "name@example.com"
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return "", "", ""
	}
	name := email[:atIndex]
	parts := strings.Split(name, ".")
	givenName := parts[0]
	familyName := parts[1]
	return givenName + " " + familyName, givenName, familyName
}

func (h *Handler) Login(c *gin.Context) {
	challenge, exists := c.GetQuery("login_challenge")
	if !exists || challenge == "" {
		c.String(http.StatusBadRequest, "Login challenge not found")
		return
	}

	loginRequest, r, err := h.hydraApi.OAuth2API.GetOAuth2LoginRequest(c).LoginChallenge(challenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.GetLoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	if loginRequest.Skip {
		fmt.Printf("Accepting login request because it was skipped\n")

		acceptLoginRequest := newAcceptLoginRequest(loginRequest.GetSubject())
		acceptResp, r, err := h.hydraApi.OAuth2API.AcceptOAuth2LoginRequest(c).LoginChallenge(challenge).AcceptOAuth2LoginRequest(*acceptLoginRequest).Execute()
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

	// Extract name from email
	name, _, _ := extractNameFromEmail(loginHint)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Challenge": challenge,
		"LoginHint": loginHint,
		"Name":      name,
	})
}

type PostLoginForm struct {
	Challenge string `form:"challenge" binding:"required"`
	Name      string `form:"name" binding:"required"`
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
		"name":           form.Name,
		"email":          form.Email,
		"email_verified": true,
	})
	acceptResp, r, err := h.hydraApi.OAuth2API.AcceptOAuth2LoginRequest(c).LoginChallenge(form.Challenge).AcceptOAuth2LoginRequest(*acceptLoginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptLoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
}
