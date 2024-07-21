package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.handleHome))
	mux.Handle("/healthz", http.HandlerFunc(app.handleHealthz))

	
	return mux
}