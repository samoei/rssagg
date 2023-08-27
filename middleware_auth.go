package main

import (
	"fmt"
	"net/http"

	"github.com/samoei/rssagg/internal/auth"
	"github.com/samoei/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Authentication Error: %v", err))
			return
		}
		user, err := cfg.DB.GetUserByAPIkey(r.Context(), apikey)

		if err != nil {
			respondWithError(w, 400, "Could not find user")
			return
		}
		handler(w, r, user)
	}
}
