package routes

import (
	hydraClient "github.com/ory/hydra-client-go"
)

type Handler struct {
	hydraApi *hydraClient.APIClient
}

func NewHandler(hydraApi *hydraClient.APIClient) *Handler {
	return &Handler{hydraApi: hydraApi}
}
