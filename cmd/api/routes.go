package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle("GET /v1/healthz", http.HandlerFunc(app.handleHealthz))
	mux.Handle("GET /v1/chirps/", http.HandlerFunc(app.handleGetChirp))
	mux.Handle("GET /v1/chirps", http.HandlerFunc(app.handleHome))
	mux.Handle("POST /v1/chirps", http.HandlerFunc(app.handleCreateChirp))
	mux.Handle("POST /v1/users", http.HandlerFunc(app.handleCreateUser))
	mux.Handle("POST /v1/login", http.HandlerFunc(app.handleLogin))
	mux.Handle("PUT /v1/users", http.HandlerFunc(app.handleUpdate))
	mux.Handle("POST /v1/refresh", http.HandlerFunc(app.handleRefresh))
	mux.Handle("POST /v1/revoke", http.HandlerFunc(app.handleRevoke))
	mux.Handle("DELETE /v1/chirps/", http.HandlerFunc(app.handleDeleteChirp))
	
	return mux
}