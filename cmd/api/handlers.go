package main

import (
	"encoding/json"
	"net/http"
	"strings"

)



func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	app.responseWithJson(w, 200, map[string]string{"message":"Hello I am from the home page"})
}


func (app *application) handleHealthz(w http.ResponseWriter, r *http.Request) {
	app.responseWithJson(w, 200, map[string]string{"message":"Yes the server is up"})
}

func (app  *application) handleCreateChirp(w http.ResponseWriter, r *http.Request){

	type reqBody struct{
		Body string `json:"body"`
	}

	params := reqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	if len(params.Body) > 140 {
		app.responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return 
	}

	words := strings.Split(params.Body, " ")
	bannedWords := []string{"kerfuffle","sharbert","fornax"}

	for i, word := range words{
		cleanedWord := word
		for _, banned := range bannedWords{
			if strings.EqualFold(strings.ToLower(word), banned){
				cleanedWord = "****"
			}
		}

		words[i] = cleanedWord 
	}

	id, err := app.chirp.InsertChirp(strings.Join(words, " "))

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "unable to create chirp")
		return 
	}

	app.responseWithJson(w, 201, map[string]int{"body":id})

}