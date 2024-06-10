package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"encoding/json"
)


//server type
type Server struct{
	fileServerHits int // stateful route
}

//healthz
func (s *Server) handleHealthz(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}


//metrics
func (s *Server) handleMetrics(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("Hits: %d", s.fileServerHits)
	res.Write([]byte(body))
}


//reset
func (s *Server) handleReset(res http.ResponseWriter, req *http.Request){
	s.fileServerHits = 0
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("File Server hit reset to 0"))
}


//app/static
func (s *Server) handleStatic(res http.ResponseWriter, req *http.Request){
	trimmedPath := strings.TrimPrefix(req.URL.Path, "/app/static/")
	http.ServeFile(res, req, "./static/" + trimmedPath)
}

//app/images
func(s *Server) handleImages(res http.ResponseWriter, req *http.Request){
	trimmedPath := strings.TrimPrefix(req.URL.Path, "/app/images")
	http.ServeFile(res, req, "./asset/" + trimmedPath)
}

//metricsMiddleware
func (s *Server) metricsMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
		s.fileServerHits++
		next.ServeHTTP(res, req)
	})
}

//logMiddleware
func (s *Server) logMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
		log.Printf("%s %s", req.Method, req.URL.Path)
		next.ServeHTTP(res, req)
	})
}

//validate chirpy
func(s *Server) handleValidateChirpy(res http.ResponseWriter, req *http.Request){
	type requestBody struct{
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := requestBody{}
	err := decoder.Decode(&params)

	if(err != nil){
		respondWithError(res, http.StatusBadRequest, "Invalid Request payload")
		return 
	}

	if len(params.Body) > 140{
		respondWithError(res, http.StatusBadRequest, "Chirpy is too long")
		return 
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedWords := params.Body 

	for _, word := range profaneWords{
		cleanedWords = strings.ReplaceAll(cleanedWords, word, "****")
		cleanedWords = strings.ReplaceAll(cleanedWords, strings.Title(word), "****")
	}

	respondWithJson(res, http.StatusOK, cleanedWords)


}

func respondWithJson(res http.ResponseWriter, code int, payload interface{}) error{
	response, err := json.Marshal(payload)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"error":"Something Went Wrong"}`))
		return err
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.WriteHeader(code)
	res.Write(response)
	return nil
}


func respondWithError(res http.ResponseWriter, code int, message string) error{
	return respondWithJson(res, code, map[string]string{"error":message})

}