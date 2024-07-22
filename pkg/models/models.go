package models 

import(
	"errors"
	// "time"
)


var ErrNoRecord = errors.New("there is no chirp found")

type Chirp struct{
    Id int;
	Body string;
}