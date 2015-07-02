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

	// Nevermind -- we have to use QueryRow to grab the user_id out of the statement, so we can make sure to
	// properly build the corresponding row in the auth schema
	query_str, err := db.Preparex(`INSERT INTO tinyplannr_api.user (email, first_name, last_name, zip_code, update_dt)
	                               VALUES ($1, $2, $3, $4, $5)
	                               RETURNING user_id;`)

	err = query_str.QueryRowx(au.Email, au.FirstName, au.LastName, au.ZipCode, au.UpdateDt).Scan(&au.UserId)

	if err != nil {
		return err
	}

	// Creates the corresponding row/entry in the auth table.
	err = CreateUserAuth(db, au)

	if err != nil {
		return err
	}

	return err
}

func DeleteUser(db *sqlx.DB, email string) error {
	var err error

	// We have to drop the UserAuth entry first to avoid breaking foreign key constraints
	err = DeleteUserAuth(db, email)

	if err != nil {
		log.Printf("Encountered an error removing the UserAuth entry: %v", err)
	}

	tx := db.MustBegin()
	tx.MustExec(`DELETE FROM tinyplannr_api.user WHERE email = $1`, email)
	err = tx.Commit()

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}


