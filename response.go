package main

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func sendResponse(w http.ResponseWriter, code int, resp *ApiResponse) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
