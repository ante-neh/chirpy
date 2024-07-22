package mysql

import (
	"database/sql"
	"github.com/ante-neh/chirpy/pkg/models"
)


type ChirpModel struct{
	Db *sql.DB;
}


func(m *ChirpModel) InsertChirp(content string) (int, error){ 
	stmt := "INSERT INTO chirpies(body) values(?)"
	result, err := m.Db.Exec(stmt, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}

	return int(id), nil
}


func(m *ChirpModel) GetChirp(id int) ( *models.Chirp, error){
	stmt := "SELECT id, body FROM chirpies WHERE id=?" 
	chirp := &models.Chirp{} 
	err := m.Db.QueryRow(stmt, id).Scan(&chirp.Id, &chirp.Body)

	if err == sql.ErrNoRows{
		return nil, models.ErrNoRecord
	}

	if err != nil{
		return nil, err
	}

	return chirp, nil 
}


func (m *ChirpModel) GetChirps()([]*models.Chirp, error){
	stmt := "SELECT * from chirpies"
	rows, err := m.Db.Query(stmt)

	if err != nil{
		return nil, err
	}

	chirps := []*models.Chirp{} 
	
	for rows.Next(){
		chirp := &models.Chirp{}
		err = rows.Scan(&chirp.Id, &chirp.Body) 

		if err != nil{
			return nil, err
		} 

		chirps = append(chirps, chirp)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}
	
	return chirps, nil 
}