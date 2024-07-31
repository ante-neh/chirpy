package models 

import(
	"errors"
	// "time"
)


var ErrNoRecord = errors.New("there is no chirp found")
var ErrNoUser = errors.New("User not found")
type Chirp struct{
    Id int;
	Body string;
}

type User struct{
	Id int
	Email string
	Password string
	RefreshToken string
}