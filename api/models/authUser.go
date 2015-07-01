package models

import (
	"time"
	// "database/sql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type AuthUser struct {
	AuthID        int64      `db:"auth_id"`
	UserID        int64      `db:"user_id"`
	Email         string     `db:"email"`
	HashPW        string     `db:"hash_pw"`
	CreateDt      time.Time  `db:"create_dt"`
	UpdateDt      time.Time  `db:"update_dt"`
	LastLoginDt   time.Time  `db:"last_login_dt"`
}

func CreateUserAuth(db *sqlx.DB, au ApiUserCreate) error {
	var err error

	// Let's use Bcrypt to generate a password hash
	pw := []byte(au.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, 13)
	fmt.Println(hashedPassword)
	if err != nil {
		return err
	}

	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO tinyplannr_auth.user (email, first_name, last_name, zip_code, update_dt)
	             VALUES ($1, $2, $3, $4, $5);`, au.Email, au.FirstName, au.LastName, au.ZipCode, au.UpdateDt)
	err = tx.Commit()

	return err
}
