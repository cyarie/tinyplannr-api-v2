package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

// Data struct for ApiUser
type ApiUser struct {
	UserId		int64 		`json:"user_id" db:"user_id"`
	Email		string		`json:"email" db:"email"`
	FirstName	string		`json:"first_name" db:"first_name"`
	LastName	string		`json:"last_name" db:"last_name"`
	ZipCode		int64		`json:"zip_code" db:"zip_code"`
	IsActive	bool		`json:"is_active" db:"is_active"`
	CreateDt	time.Time	`json:"create_dt" db:"create_dt"`
	UpdateDt	time.Time	`json:"update_dt" db:"update_dt"`
}

type ApiUserCreate struct {
	UserId		int64 		`json:"user_id" db:"user_id"`
	Email		string		`json:"email" db:"email"`
	Password    string		`json:"password" db:"password""`
	FirstName	string		`json:"first_name" db:"first_name"`
	LastName	string		`json:"last_name" db:"last_name"`
	ZipCode		int64		`json:"zip_code" db:"zip_code"`
	IsActive	bool		`json:"is_active" db:"is_active"`
	CreateDt	time.Time	`json:"create_dt" db:"create_dt"`
	UpdateDt	time.Time	`json:"update_dt" db:"update_dt"`
}

func GetUserData(db *sqlx.DB, userId int64) (*ApiUser, error) {
	var user ApiUser
	var err error

	err = db.Get(&user, "SELECT * FROM tinyplannr_api.user WHERE user_id = $1", userId)
	log.Println(err)

	return &user, err
}

func CreateUser(db *sqlx.DB, au ApiUserCreate) error {
	var err error

	// Create a transaction, and insert the JSON into the DB
	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO tinyplannr_api.user (email, first_name, last_name, zip_code, update_dt)
	             VALUES ($1, $2, $3, $4, $5);`, au.Email, au.FirstName, au.LastName, au.ZipCode, au.UpdateDt)
	err = tx.Commit()

	if err != nil {
		return err
	}

	return err
}


