package main

import (
	"fmt"
	"net/http"
)



func main() {
	mux := http.NewServeMux()
	server := &Server{}
	port := ":5000"

	//pattern [METHOD] [HOST]/[PATH] Note that all three parts are optional.



	//route
	mux.Handle("GET api/healthz", server.logMiddleware(http.HandlerFunc(server.handleHealthz)))
	mux.Handle("GET api/metrics", server.logMiddleware(http.HandlerFunc(server.handleMetrics)))
	mux.Handle("/reset", server.logMiddleware(http.HandlerFunc(server.handleReset)))
	mux.Handle("/app/static/", server.logMiddleware(server.metricsMiddleware(http.HandlerFunc(server.handleStatic))))
	mux.Handle("/app/images/", server.logMiddleware(server.metricsMiddleware(http.HandlerFunc(server.handleImages))))



	fmt.Printf("Server running on port %s \n", port)

	if err := http.ListenAndServe(port, mux); err != nil{
		fmt.Println("Server is not working")
	}
}
