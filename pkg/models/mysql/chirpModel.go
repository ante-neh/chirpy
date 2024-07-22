package mysql

import(
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


func(m *ChirpModel) GetChirp(id int) *models.Chirp{
	return nil
}