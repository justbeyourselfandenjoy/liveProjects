package main

import (
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthHandler has been called")
	w.WriteHeader(http.StatusOK)
	return
}
