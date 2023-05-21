package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func SafeJsonEncode(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		_, err := w.Write([]byte("response JSON encoding error"))
		if err != nil {
			log.Printf("Encoding JSON crictical error: %v", err)
		}
	}
}
