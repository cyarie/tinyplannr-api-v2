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

func GetUserData(db *sqlx.DB, userId int64) (*ApiUser, error) {
	var user ApiUser
	var err error

	err = db.Get(&user, "SELECT * FROM tinyplannr_api.user WHERE user_id = $1", userId)
	log.Println(err)

	return &user, err
}


