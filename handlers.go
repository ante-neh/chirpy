package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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