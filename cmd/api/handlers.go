package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/ante-neh/chirpy/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
		return 
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


func(app *application) handleCreateUser(w http.ResponseWriter, r *http.Request){
	type reqeustBody struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body) 
	params := reqeustBody{} 
	err := decoder.Decode(&params) 

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost) 
	if err != nil{
		app.responseWithError(w, 500, "unable to hash the password")
		return 
	}


	lastId, err := app.chirp.CreateUser(params.Email, string(hashedPassword))

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "Unable to create a user")
		return 
	}

	app.responseWithJson(w, 201, map[string]int{"UserId":lastId})
	
}


func (app *application) handleLogin(w http.ResponseWriter, r *http.Request){
	type reqBody struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		Expires_in_seconds  *int   `json:"expires_in_seconds"`
	}

	params := reqBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}


	user, err := app.chirp.UserLogin(params.Email)
	if err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	

	if err := bcrypt.CompareHashAndPassword([]byte(user.Email), []byte(params.Password)); err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, 401, "Unauthorized user")
		return 
	}


	var expires int 

	if params.Expires_in_seconds != nil{
		if *params.Expires_in_seconds > 86000{
			expires = 86000
		}
	}else{
		expires = 86000
	}


	token, err := app.createJWT(strconv.Itoa(user.Id), expires)

	if err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, 500, "Somthing went wrong")
		return 
	}

	type response struct{
		id int 
		email string
		token string
	}

	res := response{
		id:user.Id,
		email:user.Email,
		token:token,
	}
	

	app.responseWithJson(w, 200, res)
	
}



func(app *application) handleUpdate(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization") 
	token := strings.TrimPrefix(authorization, "Bearer")
	app.infoLog.Println(token)

}


func (app *application) createJWT(id string, expires_at int)(string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt:jwt.NewNumericDate(time.Now().Add(time.Duration(expires_at) * time.Second)), 
		Issuer:"Chirpy",
		IssuedAt:jwt.NewNumericDate(time.Now()),
		Subject: id,
		}

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.jwt)) 

	if err != nil{
		return "", err
	}
	
	return tokenString, nil

}