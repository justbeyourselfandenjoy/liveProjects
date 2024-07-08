package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthHandler has been called")
	w.WriteHeader(http.StatusOK)
	return
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("orderHandler has been called")

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content-type 'application/json' is allowed, got '%s'", contentType)
		return
	}

	var orderReceived Order

	err = json.Unmarshal(bodyBytes, &orderReceived)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("%v\n", orderReceived)
	return
}
