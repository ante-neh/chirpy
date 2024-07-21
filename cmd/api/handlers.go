package main

import "net/http"



func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.responseWithJson(w, 200, map[string]string{"message":"Hello I am from the home page"})
}


func (app *application) handleHealthz(w http.ResponseWriter, r *http.Request) {
	app.responseWithJson(w, 200, map[string]string{"message":"Yes the server is up"})
}
