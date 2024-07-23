package model

import (
	"database/sql"
	"errors"
	"os"
	_ "github.com/jackc/pgx/v5/stdlib"
	sq "github.com/Masterminds/squirrel"
)

type User struct {
	Id         int    `json:"user_id" form:"user_id"`
	Username   string `json:"user_name" form:"user_name"`
	FirstName  string `json:"first_name" form:"first_name"`
	LastName   string `json:"last_name" form:"last_name"`
	Email      string `json:"email" form:"email"`
	Status     string `json:"user_status" form:"user_status"`
	Department string `json:"department" form:"department"`
}

type UserNoId struct {
	Username   string `json:"user_name" form:"user_name"`
	FirstName  string `json:"first_name" form:"first_name"`
	LastName   string `json:"last_name" form:"last_name"`
	Email      string `json:"email" form:"email"`
	Status     string `json:"user_status" form:"user_status"`
	Department string `json:"department" form:"department"`
}

func openDb() (*sql.DB, error) {
	url, found := os.LookupEnv("DATABASE_URL")
	if !found {
		url = "postgres://admin:password@localhost/GoAngularCrud"
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetUsers() ([]User, error) {
	db, dberr := openDb()
	if dberr != nil {
		return nil, dberr
	}
	defer db.Close()
	rows, sqlerr := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("user_id", "user_name", "first_name", "last_name", "email", "user_status", "department").
			From("users").
		RunWith(db).
		Query()
	if sqlerr != nil {
		return nil, sqlerr
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Department)
		users = append(users, user)
	}
	return users, nil
}

func GetUser(id int) (*User, error) {
	db, dberr := openDb()
	if dberr != nil {
		return nil, dberr
	}
	defer db.Close()
	var user User
	sqlerr := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("user_id", "user_name", "first_name", "last_name", "email", "user_status", "department").
			From("users").
			Where(sq.Eq{"user_id": id}).
			Limit(1).
		RunWith(db).
		QueryRow().
		Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Department)
	if sqlerr != nil {
		if errors.Is(sqlerr, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, sqlerr
	}
	return &user, nil
}

func (user UserNoId) Post() (*User, error) {
	db, dberr := openDb()
	if dberr != nil {
		return nil, dberr
	}
	defer db.Close()
	var id int = 0
	err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("users").
			Columns("user_name",   "first_name",   "last_name",   "email",    "user_status", "department").
			Values( user.Username, user.FirstName, user.LastName, user.Email, user.Status,   user.Department).
			Suffix("RETURNING \"user_id\"").
		RunWith(db).
		QueryRow().
		Scan(&id)
	if err != nil {
		return nil, err
	}
	return &User{
		Id: id,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Status: user.Status,
		Department: user.Department,
	}, nil
}

func (user UserNoId) Put(id int) (bool, error) {
	db, dberr := openDb()
	if dberr != nil {
		return false, dberr
	}
	defer db.Close()
	var inserted bool
	err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("users").
			Columns("user_id", "user_name",   "first_name",   "last_name",   "email",    "user_status", "department").
			Values( id,        user.Username, user.FirstName, user.LastName, user.Email, user.Status,   user.Department).
			Suffix("ON CONFLICT(user_id) DO UPDATE SET user_name = EXCLUDED.user_name, first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name, email = EXCLUDED.email, user_status = EXCLUDED.user_status, department = EXCLUDED.department RETURNING (xmax = 0) AS inserted").
		RunWith(db).
		QueryRow().
		Scan(&inserted)
	if err != nil {
		return false, err
	}
	return inserted, nil
}

func DeleteUser(id int) (bool, error) {
	db, dberr := openDb()
	if dberr != nil {
		return false, dberr
	}
	defer db.Close()
	result, sqlerr := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("users").
			Where(sq.Eq{"user_id": id}).
			//Limit(1). //postgres doesn't allow limit in delete, enable this for mysql and others
		RunWith(db).
		Exec()
	if sqlerr != nil {
		return false, sqlerr
	}
	affected, afferr := result.RowsAffected()
	if afferr != nil {
		return true, nil
	}
	return affected > 0, nil
}