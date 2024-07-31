package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	chirps, err := app.chirp.GetChirps() 
	if err != nil{
		app.responseWithError(w, 500, "Something went wrong")
	}
	app.responseWithJson(w, 200, map[string][]*models.Chirp{"chirps": chirps})
 
}



func (app *application) handleCreateChirp(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("authorization")
	tokenString := strings.TrimPrefix(authorization, "Bearer")
	claims, er := app.validateToken(tokenString) 

	if er != nil{
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	type reqBody struct{
		Body string `json:"body"`
	}

	params := reqBody{}
	err := json.NewDecoder(r.Body).Decode(&params)

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

	id, err := app.chirp.InsertChirp(strings.Join(words, " "), claims.Subject)

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
			return 
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

	
	params := reqeustBody{} 
	err := json.NewDecoder(r.Body).Decode(&params) 

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost) 
	if err != nil{
		app.responseWithError(w, 500, "unable to hash the password")
		return 
	}

	refreshToken, err := app.generateRefreshToken(5)
	if err != nil{
		app.responseWithError(w, 500, "unable to generate refresh Token")
		return 
	}


	lastId, err := app.chirp.CreateUser(params.Email, string(hashedPassword), refreshToken)

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

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}


	user, err := app.chirp.UserLogin(params.Email)
	if err == models.ErrNoUser{
		app.responseWithError(w, 404, "Account Not Exist")
		return 
	}

	if err != nil{
		app.errorLog.Println(err)
		app.responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return 
	}

	

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil{
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
		app.responseWithError(w, 500, "Somthing went wrong")
		app.errorLog.Fatal(err)
		return 
	}

	type response struct{
		Id int 
		Email string
		Token string
		RefreshToken string 
	}

	res := &response{
		Id:user.Id,
		Email:params.Email,
		Token:token,
		RefreshToken:user.RefreshToken,
	}
	

	app.responseWithJson(w, 200, res)
	
}



func(app *application) handleUpdate(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization") 
	token := strings.TrimPrefix(authorization, "Bearer")
	claims, err := app.validateToken(token)

	if err != nil{
		app.responseWithError(w, http.StatusInternalServerError, "invalid token")
	}

	id, er := app.chirp.UpdateChirp(claims.Subject)

	if er != nil{
		app.responseWithError(w, 400, "Unathorized")
		return 
	}


	app.responseWithJson(w, 200, map[string]int{"id":id})

}


func (app *application) handleDeleteChirp(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("authorization")
	tokenString := strings.TrimPrefix(authorization, "Bearer")
	claims, err := app.validateToken(tokenString)
	if err != nil{
		app.responseWithError(w, 403, "Unauthorized")
		return 
	}

	err = app.chirp.DeleteChirp(claims.Subject)
	if err != nil{
		return 
	}

	app.responseWithJson(w, 200, map[string]string{"message":"chirp deleted"})
}

func (app *application) handleRefresh(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization")
	refresh := strings.TrimPrefix(authorization, "Bearer") 
	user, err := app.chirp.GetRefreshToken(refresh)

	if err != sql.ErrNoRows{
		app.responseWithError(w, 404, "user doen't exist")
		return 
	}
	if err != nil{
		app.responseWithJson(w, 500, "Internal server error")
		return 
	}

	token, err := app.createJWT(strconv.Itoa(user.Id), 3600 ) 

	if err != nil{
		app.responseWithError(w, 500, "Unable to create JWT")
		return 

	}
    app.responseWithJson(w, 200, map[string]string{"token":token})

}

func (app *application) handleRevoke(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization")
	refreshToken := strings.TrimPrefix(authorization, "Bearer")
	err := app.chirp.RevokeToken(refreshToken)

	if err != nil{
		app.responseWithError(w, 500, "Internal Server Error")
	}

	app.responseWithJson(w, 204, "") 

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

func (app *application) generateRefreshToken(length int) (string, error){
	bytes := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, bytes)

	if err != nil{
		return "", nil
	}

	return hex.EncodeToString(bytes), nil 
}


func (app *application) validateToken(tokenString string) (*jwt.RegisteredClaims, error){
	claims := &jwt.RegisteredClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(interface{},error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return app.jwt, nil 
	})

	if err != nil{
		return nil, err
	}

	if !token.Valid{
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}