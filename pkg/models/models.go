package models 

import(
	"errors"
	// "time"
)


var ErrNoRecord = errors.New("")

type Chirp struct{
	id int;
	content string;
}