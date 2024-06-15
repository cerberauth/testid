package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go/v2"
)

func newAcceptConsentRequest(consentRequest *hydraClient.OAuth2ConsentRequest) *hydraClient.AcceptOAuth2ConsentRequest {
	acceptConsentRequest := hydraClient.NewAcceptOAuth2ConsentRequest()
	acceptConsentRequest.SetRemember(true)
	acceptConsentRequest.SetRememberFor(3600 * 12)
	acceptConsentRequest.SetGrantScope(consentRequest.GetRequestedScope())
	acceptConsentRequest.SetGrantAccessTokenAudience(consentRequest.GetRequestedAccessTokenAudience())
	return acceptConsentRequest
}

func (h *Handler) Consent(c *gin.Context) {
	challenge, exists := c.GetQuery("consent_challenge")
	if !exists || challenge == "" {
		c.String(http.StatusBadRequest, "Consent challenge not found")
		return
	}

	consentRequest, r, err := h.hydraApi.OAuth2API.GetOAuth2ConsentRequest(c).ConsentChallenge(challenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	if consentRequest.GetSkip() {
		fmt.Printf("Accepting consent request because it was skipped\n")

		acceptConsentRequest := newAcceptConsentRequest(consentRequest)
		acceptResp, r, err := h.hydraApi.OAuth2API.AcceptOAuth2ConsentRequest(c).ConsentChallenge(challenge).AcceptOAuth2ConsentRequest(*acceptConsentRequest).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptConsentRequest``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
	}

	logoUri := consentRequest.GetClient().LogoUri
	if *logoUri == "" {
		logoUri = nil
	}
	clientName := consentRequest.GetClient().ClientName
	tosUri := consentRequest.GetClient().TosUri
	if *tosUri == "" {
		tosUri = nil
	}
	policyUri := consentRequest.GetClient().PolicyUri
	if *policyUri == "" {
		policyUri = nil
	}

	c.HTML(http.StatusOK, "consent.html", gin.H{
		"Challenge":      challenge,
		"RequestedScope": consentRequest.GetRequestedScope(),
		"LogoUri":        logoUri,
		"ClientName":     clientName,
		"TosUri":         tosUri,
		"PolicyUri":      policyUri,
	})
}

type PostConsentForm struct {
	Challenge string `form:"challenge" binding:"required"`
	// Scopes    []string `form:"scopes" binding:"required"`
}

func (h *Handler) PostConsent(c *gin.Context) {
	var form PostConsentForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	consentRequest, r, err := h.hydraApi.OAuth2API.GetOAuth2ConsentRequest(c).ConsentChallenge(form.Challenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.GetConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	acceptConsentRequest := newAcceptConsentRequest(consentRequest)
	acceptResp, r, err := h.hydraApi.OAuth2API.AcceptOAuth2ConsentRequest(c).ConsentChallenge(form.Challenge).AcceptOAuth2ConsentRequest(*acceptConsentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.AcceptConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, acceptResp.GetRedirectTo())
}
