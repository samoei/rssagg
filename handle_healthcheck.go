package main

import "net/http"

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	repondWithJSON(w, 200, struct{}{})
}
