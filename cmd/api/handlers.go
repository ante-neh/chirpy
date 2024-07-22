package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ante-neh/chirpy/pkg/models"
)



func (app *application) handleHealthz(w http.ResponseWriter, r *http.Request) {
	app.responseWithJson(w, 200, map[string]string{"message":"Yes the server is up"})
}


func (app *application) handleHome(w http.ResponseWriter, r *http.Request) {
	chirps, err := app.chirp.GetChirps () 
	if err != nil{
		app.responseWithError(w, 500, "Something went wrong")
	}
	app.responseWithJson(w, 200, map[string][]*models.Chirp{"chirps": chirps})
}



func (app *application) handleCreateChirp(w http.ResponseWriter, r *http.Request){

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


func (app *application) handleGetChirp(w http.ResponseWriter, r *http.Request){
	url := strings.TrimPrefix(r.URL.Path, "/chirps/")
	id, err := strconv.Atoi(url)

	if err != nil{
		app.responseWithError(w, 404, "bad request")
	}

	result, err := app.chirp.GetChirp(id) 

	if err != nil {
		if err == models.ErrNoRecord{

			app.responseWithError(w, 404, "Chirp not found")
		}

		app.responseWithError(w, 500, "Something went wrong")
		return 
	}

	app.responseWithJson(w, 200, result)
}