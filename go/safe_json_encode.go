package main

import (
	"encoding/json"
	"net/http"
)

func SafeJsonEncode(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		_, err = w.Write([]byte("response JSON encoding error"))
	}
}
