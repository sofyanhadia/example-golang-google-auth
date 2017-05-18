package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/errors"
	uuid "github.com/satori/go.uuid"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/structs"
)

var (
	url string
	db  *sql.DB
	dot *dotsql.DotSql
)

const (
	userDBSchema     = "./database/schema.sql"
	insertQuery      = "insert"
	selectEmailQuery = "select-email"
)

func New(url string) {
	url = url
	var d *sql.DB
	var err error

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		panic("Invalid database url")
	}

	if len(parts[1]) == 0 {
		panic("Invalid database name")
	}

	if dot, err = dotsql.LoadFromFile(userDBSchema); err != nil {
		log.Fatal(err)
	}

	if d, err = sql.Open("mysql", url); err != nil {
		log.Fatal(err)
	}

	db = d
}

func Read(email string) (*structs.User, error) {
	user := &structs.User{}

	row, err := dot.QueryRow(db, selectEmailQuery, email)

	// Scan => take data
	if err := row.Scan(&user.Id, &user.Sub, &user.Name, &user.GivenName, &user.FamilyName, &user.Profile, &user.Picture, &user.Email, &user.EmailVerified, &user.Gender); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound(email, err.Error())
		}

		return nil, errors.InternalServerError(email, err.Error())
	}

	return user, err
}

func Create(user *structs.User) (sql.Result, error) {
	if _, err := Read(user.Email); err == nil {
		return nil, fmt.Errorf("User already exists!")
	}

	user.Id = uuid.NewV4().String()
	result, err := dot.Exec(db, insertQuery, user.Id, user.Sub, user.Name, user.GivenName, user.FamilyName, user.Profile, user.Picture, user.Email, user.EmailVerified, user.Gender)

	if err != nil {
		return nil, errors.InternalServerError("", err.Error())
	}

	return result, err
}