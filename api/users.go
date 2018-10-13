package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	m "proto-game-server/models"

	validate "github.com/asaskevich/govalidator"
	"github.com/badoux/checkmail"
)

type IUserStorage interface {
	Add(user *m.User) *ApiResponse

	Remove(user *m.User) *ApiResponse

	Update(user *m.User) *ApiResponse

	Get(slug string) *ApiResponse
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db}
}

func ScanUserFromRow(row *sql.Row) (*m.User, error) {
	user := new(m.User)
	err := row.Scan(&user.Id, &user.Nickname, &user.Password, &user.Fullname, &user.Email)

	return user, err
}

// nice func to remove repeating code
func throwError(code int, message string) *ApiResponse {
	log.Println(message)
	return &ApiResponse{
		Code: http.StatusBadRequest,
		Response: &m.Error{
			Code:    http.StatusBadRequest,
			Message: message}}
}

func validateUser(user *m.User) (err error) {
	// this defer catches panics from smtp module
	defer func() error {
		if rec := recover(); rec != nil {
			switch x := rec.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown error")
			}
		}
		return err
	}()

	// user fields validation
	_, err = validate.ValidateStruct(user)
	if err != nil {
		return err
	}

	// check if the email is resolvable
	err = checkmail.ValidateHost(user.Email)
	if err != nil {
		return err
	}

	return nil
}

// TODO: user validation
func (u *UserStorage) Add(user *m.User) *ApiResponse {
	if err := validateUser(user); err != nil {
		return throwError(http.StatusBadRequest, err.Error())
	}

	result, err := u.db.Exec(
		"INSERT INTO user(nickname, password, email, fullname) VALUES ($1,$2,$3,$4);",
		user.Nickname, user.Password, user.Email, user.Fullname)

	if err != nil {
		return throwError(http.StatusConflict, err.Error())
	}

	user.Id, _ = result.LastInsertId()
	return &ApiResponse{Code: 201, Response: user}
}

// TODO: user remove
func (u *UserStorage) Remove(user *m.User) *ApiResponse {
	return &ApiResponse{Code: 400, Response: &m.Error{1, "unimplemented api"}}
}

//untested. Скорее всего не работает
// TODO: user update() validation
func (u *UserStorage) Update(user *m.User) *ApiResponse {
	row := u.db.QueryRow("SELECT id, nickname, password, fullname, email FROM user WHERE id=$1", user.Id)
	oldUser, err := ScanUserFromRow(row)

	if err != nil {
		log.Println(err.Error())
		return &ApiResponse{
			Code:     http.StatusNotFound,
			Response: &m.Error{Code: http.StatusNotFound, Message: err.Error()},
		}
	}

	if user.Nickname == "" {
		user.Nickname = oldUser.Nickname
	}

	if user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}

	if user.Password == "" {
		user.Password = oldUser.Password
	}

	if user.Email == "" {
		user.Email = oldUser.Email
	}

	_, err = u.db.Exec("UPDATE user SET nickname=$1, fullname=$2, password=$3, email=$4 WHERE id=$5",
		user.Nickname, user.Fullname, user.Password, user.Email, user.Id)
	if err != nil {
		log.Println(err.Error())
		return &ApiResponse{
			Code:     http.StatusConflict,
			Response: &m.Error{Code: http.StatusConflict, Message: err.Error()},
		}
	}

	return &ApiResponse{
		Code:     http.StatusOK,
		Response: user,
	}
}

// TODO: method for recieving user's info
func (u *UserStorage) Get(slug string) *ApiResponse {
	return &ApiResponse{Code: 400, Response: &m.Error{1, "unimplemented api"}}
}

// UPDATE user SET email = $1, WHERE nickname = $2
// curl -X PUT -d '{"nickname":"asd1, "email":"kek@kek.os"}' localhost:8080/user/
