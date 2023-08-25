package main

import "net/http"

func handleErr(w http.ResponseWriter, r *http.Request) {
	repondWithJSON(w, 400, "Something went wrong")
}
