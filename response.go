package main

import (
	"encoding/json"
	"net/http"
)

// JSONError write a error to ResponseWriter
func JSONError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// JSON Write to ResponseWriter the res interface
func JSON(w http.ResponseWriter, res interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
