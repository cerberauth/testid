package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cerberauth/testid/hydra-login-consent/routes"
	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go/v2"
)

func setupHydraClient() *hydraClient.APIClient {
	hydraAdminURL := os.Getenv("HYDRA_ADMIN_URL")
	if hydraAdminURL == "" {
		hydraAdminURL = "http://localhost:4445"
	}

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

func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/css/") {
			c.Header("Cache-Control", "public, max-age=86400, s-maxage=86400, stale-while-revalidate=86400, stale-if-error=86400")
		} else if c.Request.URL.Path == "/" {
			c.Header("Cache-Control", "public, s-maxage=86400, stale-while-revalidate=86400, stale-if-error=86400")
		} else {
			c.Header("Cache-Control", "no-store, max-age=0")
		}
		c.Next()
	}
}

func setupRouter(h *routes.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(cacheMiddleware())

	r.GET("/", h.Index)
	r.GET("/error", h.Error)
	r.GET("/login", h.Login)
	r.POST("/login", h.PostLogin)
	r.GET("/consent", h.Consent)
	r.POST("/consent", h.PostConsent)
	r.GET("/logout", h.Logout)

	r.Static("/css", "./static/css")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.StaticFile("/robots.txt", "./static/robots.txt")

	return r
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
