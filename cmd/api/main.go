package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ante-neh/chirpy/pkg/models/mysql"
	"github.com/joho/godotenv"
)

type application struct{
	infoLog *log.Logger 
	errorLog *log.Logger
	chirp *mysql.ChirpModel
	jwt string
}

func main() {
	
	infoLog := log.New(os.Stdout, "Info\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate | log.Ltime | log.Lshortfile)
	
	
	err := godotenv.Load()
	if err != nil{
		errorLog.Fatal("Unable to extract .env file")
	}
	
	
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DNS:= os.Getenv("DNS")
	ADDRESS := os.Getenv("ADDRESS")



	address := flag.String("address", ADDRESS, "Port number where the server is accessible")
	dns := flag.String("dns", DNS,"connection string")
	flag.Parse()



	db, e := openDb(*dns)

	if e != nil{
		errorLog.Fatal(e)
	}

	defer db.Close() 


	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
		chirp: &mysql.ChirpModel{Db:db},
		jwt:JWT_SECRET,
	}


	

	server := &http.Server{
		Addr: *address,
		Handler : app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Server is running on port: %v", *address)
	err = server.ListenAndServe() 
	errorLog.Fatal(err)
}                                                                                          


func openDb(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil{
		
		return nil, err
	}

	if err := db.Ping(); err !=nil {
		return nil, err
	}

	return db, nil
}