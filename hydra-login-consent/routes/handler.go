package routes

import (
	hydraClient "github.com/ory/hydra-client-go/v2"
)

type Handler struct {
	hydraApi *hydraClient.APIClient
}

func NewHandler(hydraApi *hydraClient.APIClient) *Handler {
	return &Handler{hydraApi: hydraApi}
}
