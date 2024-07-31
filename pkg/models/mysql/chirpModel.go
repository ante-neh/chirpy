package mysql

import (
	"database/sql"
	"fmt"
	"github.com/ante-neh/chirpy/pkg/models"
)


type ChirpModel struct{
	Db *sql.DB;
}


func (m *ChirpModel) CreateUser(email,refreshToken,  password string) (int, error){
	stmt := "INSERT INTO user(email,password, refreshToken) VALUES(?, ?, ?)"
	result, err := m.Db.Exec(stmt, email, password, refreshToken)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId() 
	if err != nil{
		return 0, err
	}

	return int(lastId), nil
}

func (m *ChirpModel) UserLogin(email string) (*models.User, error){
	stmt := "SELECT id, password, email, refreshToken FROM user WHERE email=?"
	u := &models.User{}
	err := m.Db.QueryRow(stmt, email).Scan(&u.Id, &u.Email, &u.Password, &u.RefreshToken)
	
	if err == sql.ErrNoRows{
		fmt.Println(err)
		return nil, models.ErrNoUser
	}

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	return u, nil 
}

func(m *ChirpModel) InsertChirp(content string, id string) (int, error){ 
	stmt := "INSERT INTO chirpies(body, userId) values(?, ?)"
	result, err := m.Db.Exec(stmt, content, id)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}

	return int(lastId), nil
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

	defer rows.Close() 

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


func (m *ChirpModel) DeleteChirp(id string)(error){
	stmt := "DELETE FROM chirps WHERE userId = ?"
	_, err := m.Db.Exec(stmt, id) 
	if err != nil{
		return err
	}

	return nil 
}


func (m *ChirpModel) UpdateChirp(id string)(int, error){
	return 0, nil
}
func (m *ChirpModel) GetRefreshToken(refreshToken string) (*models.User, error){
	stmt := "SELECT refreshToken, id, email, password FROM user WHERE refreshToken = ?"
	u := &models.User{}
	err := m.Db.QueryRow(stmt, refreshToken).Scan(&u.RefreshToken, &u.Id ,&u.Email, &u.Password)
	
	if err == sql.ErrNoRows{
		return nil, models.ErrNoUser
	}

	if err != nil{
		return nil, err
	}


	return u, nil
}


func (m *ChirpModel) RevokeToken(refreshToken string) error {
	stmt := "DELETE FROM user WHERE refreshToken = ?"
	_, err := m.Db.Exec(stmt, refreshToken)

	if err != nil{
		return err 
	}

	return nil 
}
