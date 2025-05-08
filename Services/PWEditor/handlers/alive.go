package handlers

import (
	"libery-dungeon-libs/libs/libery_networking"
	"net/http"
)

func AliveHandler(server_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		return
	}
}
