package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle("GET /v1/healthz", http.HandlerFunc(app.handleHealthz))
	mux.Handle("GET /v1/chirps", http.HandlerFunc(app.handleHome))
	mux.Handle("POST /v1/chirps", http.HandlerFunc(app.handleCreateChirp))
	mux.Handle("GET /v1//chirps/", http.HandlerFunc(app.handleGetChirp))
	mux.Handle("POST /v1/users", http.HandlerFunc(app.handleCreateUser))
	return mux
}