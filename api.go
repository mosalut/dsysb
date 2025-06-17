package main

import (
	"net/http"
	"encoding/json"
)

const API_NOT_FOUND = "API not found"

type responseResult_T struct {
	Success bool	`json:"success"`
	Message string	`json:"message"`
	Data	[]byte	`json:"data"`
}

func writeResult(w http.ResponseWriter, result responseResult_T) {
	rr, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(rr)
}

type responseResult2_T struct {
	Success bool	`json:"success"`
	Message string	`json:"message"`
	Data	interface{}	`json:"data"`
}

func writeResult2(w http.ResponseWriter, result responseResult2_T) {
	rr, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(rr)
}
