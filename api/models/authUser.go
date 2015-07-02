package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthUser struct {
	AuthID      int64     `db:"auth_id"`
	UserID      int64     `db:"user_id"`
	Email       string    `db:"email"`
	HashPW      string    `db:"hash_pw"`
	CreateDt    time.Time `db:"create_dt"`
	UpdateDt    time.Time `db:"update_dt"`
	LastLoginDt time.Time `db:"last_login_dt"`
}

func CreateUserAuth(db *sqlx.DB, au ApiUserCreate) error {
	var err error

	// Let's use Bcrypt to generate a password hash
	pw := []byte(au.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, 13)

	if err != nil {
		return err
	}

	// Rolling a transaction here. We don't need to return anything.
	tx := db.MustBegin()
	tx.MustExec(`INSERT INTO tinyplannr_auth.user (user_id, email, hash_pw, update_dt, last_login_dt)
	             VALUES ($1, $2, $3, $4, $5);`, au.UserId, au.Email, string(hashedPassword), time.Now(), time.Now())
	err = tx.Commit()

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func DeleteUserAuth(db *sqlx.DB, email string) error {
	var err error

	tx := db.MustBegin()
	tx.MustExec(`DELETE FROM tinyplannr_auth.user WHERE email = $1`, email)
	err = tx.Commit()

	if err != nil {
		log.Printf("Error removing user from auth table: %v", err)
		return err
	}

	return err
}
