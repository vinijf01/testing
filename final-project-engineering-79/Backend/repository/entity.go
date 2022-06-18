package repository

import "time"

// sama dengan field pada db
//untuk menampung data pada golang butuh sebuah struct

type User struct {
	ID        int
	Username  string
	Email     string
	Nohp      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
