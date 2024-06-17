package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Chirp struct {
	ID int 
	Body string
}


type DB struct {
	path string
	mux *sync.RWMutex
}

type DbStructure struct{
	Chirps map[int]Chirp 
}

func newDb(path string) (*DB, error){
	db := &DB{
		path:path,
		mux : &sync.RWMutex{},
	}

	if err := db.ensureDb(); err != nil{
		return nil, err
	}

	return db, nil
}

func (db *DB) ensureDb() error{
	db.mux.Lock()
	defer db.mux.Unlock()

	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist){
		emptyStructure := DbStructure{
			Chirps:make(map[int]Chirp),
		}

		return db.WriteDb(emptyStructure)
	}
	return nil
}

func (db *DB) WriteDb(dbStructure DbStructure) error{
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil{
		return err
	}


	os.WriteFile(db.path, data, 0644)

	return nil
}

func (db *DB) loadDb() (DbStructure, error){
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil{
		return DbStructure{}, err
	}

	var dbStructure DbStructure

	 if err := json.Unmarshal(data, &dbStructure); err != nil {
		return DbStructure{}, nil
	 }

	 return dbStructure, nil 
}