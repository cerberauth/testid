package main

import (
	"fmt"
	"os"

	"github.com/cerberauth/testid/hydra-login-consent/routes"
	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go"
)

func setupHydraClient() *hydraClient.APIClient {
	hydraAdminURL := os.Getenv("HYDRA_ADMIN_URL")
	if hydraAdminURL == "" {
		hydraAdminURL = "http://localhost:4445"
	}
	fmt.Printf("Hydra Admin URL: %s\n", hydraAdminURL)

	configuration := hydraClient.NewConfiguration()
	configuration.Debug = hydraAdminURL == "http://localhost:4445"
	configuration.Servers = []hydraClient.ServerConfiguration{
		{
			URL: hydraAdminURL,
		},
	}

	hydraAdminClient := hydraClient.NewAPIClient(configuration)
	return hydraAdminClient
}

func setupRouter(h *routes.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/login", h.Login)
	r.POST("/login", h.PostLogin)
	r.GET("/consent", h.Consent)
	r.POST("/consent", h.PostConsent)
	r.GET("/logout", h.Logout)

	return r
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	hydraAdminClient := setupHydraClient()
	h := routes.NewHandler(hydraAdminClient)
	r := setupRouter(h)
	r.LoadHTMLGlob("templates/*")

	fmt.Printf("Server listening on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "Error when starting the server: %v\n", err)
	}
}
