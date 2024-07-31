package main

import (
	"encoding/json"
	"net/http"
)

func(app *application) responseWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([] byte("Something went wrong"))
		return err
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	
	return nil
}



func (app *application) responseWithError(w http.ResponseWriter, code int, message string) error {
	return app.responseWithJson(w, code, map[string]string{"error":message})
}