package repository

//untuk mengakses ke db

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	FetchUsers() ([]User, error)
	InsertUser(userRequest UserRequest) (User, error) //untuk insert data user ke db
	LoginUser(email string, password string) (*string, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) FetchUsers() ([]User, error) {
	var users []User

	row, err := u.db.Query("SELECT id, username, email, nohp, password FROM users")
	if err != nil {
		return users, err
	}
	for row.Next() {
		var user User

		err := row.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Nohp,
			&user.Password,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepository) InsertUser(userRequest UserRequest) (User, error) {
	sqlStatement := `INSERT INTO users (username, email, nohp, password) VALUES (?,?,?,?)`

	res, err := u.db.Prepare(sqlStatement)

	if err != nil {
		return User{}, err
	}
	defer res.Close()

	newRes, err := res.Exec(
		userRequest.Username,
		userRequest.Email,
		userRequest.Nohp,
		userRequest.Password,
	)
	fmt.Println("succes", newRes)
	newUser := User{
		Username: userRequest.Username,
		Email:    userRequest.Email,
		Nohp:     userRequest.Nohp,
		Password: userRequest.Password,
	}
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (u *UserRepository) LoginUser(email string, password string) (*string, error) {
	users, err := u.FetchUsers()

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			return &user.Email, nil
		}
	}
	return nil, fmt.Errorf("Login Failed")
}
