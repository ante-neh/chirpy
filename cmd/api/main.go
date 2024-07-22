package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ante-neh/chirpy/pkg/models/mysql"
)

type application struct{
	infoLog *log.Logger 
	errorLog *log.Logger
	chirp *mysql.ChirpModel
}

func main() {
	
	address := flag.String("address", ":4000", "Port number where the server is accessible")
	dns := flag.String("dns", "anteneh:1919@/chirpy?parseTime=True","connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "Info\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate | log.Ltime | log.Lshortfile)

	db, e := openDb(*dns)

	if e != nil{
		errorLog.Fatal(e)
	}

	defer db.Close() 


	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
		chirp: &mysql.ChirpModel{Db:db},
	}


	

	server := &http.Server{
		Addr: *address,
		Handler : app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Server is running on port: %v", *address)
	err := server.ListenAndServe() 
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