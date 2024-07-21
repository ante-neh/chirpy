package mysql

import(
	"database/sql"
	"github.com/ante-neh/pkg/models"
)


type ChirpModel struct{
	Db *sql.DB;
}


func(m *ChirpModel) InserChirp(content string) (int, error){
	return 1, nil
}