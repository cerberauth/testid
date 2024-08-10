package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/peterhellberg/link"

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

func cleanupClients(ctx context.Context, hydraAdminClient *hydraClient.APIClient) (int, error) {
	deletedClients := 0
	pageSize := 200
	pageToken := "1"

	fmt.Println("Cleaning up clients...")
	for {
		clients, r, err := hydraAdminClient.OAuth2API.ListOAuth2Clients(ctx).PageSize(int64(pageSize)).PageToken(pageToken).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `OAuth2Api.ListOAuth2Clients``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

			return deletedClients, err
		}

		for _, client := range clients {
			if client.GetMetadata() != nil {
				metadata, ok := client.GetMetadata().(map[string]interface{})
				if ok && metadata["disable_cleanup"] == "true" {
					continue
				}
			}

			// Check if client creation time is older than 1 day
			if time.Since(client.GetCreatedAt()).Hours() > 24 {
				_, err := hydraAdminClient.OAuth2API.DeleteOAuth2Client(ctx, client.GetClientId()).Execute()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting client: %v\n", err)
					continue
				}
				deletedClients++
			}
		}

		next := link.ParseHeader(r.Header)["next"]
		if next == nil {
			break
		}

		nextUri, err := url.Parse(next.URI)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing next URI: %v\n", err)
			return deletedClients, err
		}
		pageToken = nextUri.Query().Get("page_token")
	}

	return deletedClients, nil
}

func main() {
	ctx := context.Background()
	hydraAdminClient := setupHydraClient()

	deletedClients, err := cleanupClients(ctx, hydraAdminClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning up clients: %v\n", err)
		return
	}

	fmt.Printf("Deleted %d clients\n", deletedClients)
}
