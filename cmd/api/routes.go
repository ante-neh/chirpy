package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(app.handleHome))
	mux.Handle("GET /healthz", http.HandlerFunc(app.handleHealthz))
	mux.Handle("POST /chirps", http.HandlerFunc(app.handleCreateChirp))
	
	return mux
}