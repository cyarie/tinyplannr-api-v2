package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Data struct for ApiEvent
type ApiEvent struct {
	Email        string    `json:"email"`
	Title        string	   `json:"title" db:"title"`
	Description  string	   `json:"description" db:"description"`
	Location     string    `json:"location" db:"location"`
	AllDay       bool      `json:"all_day" db:"all_day"`
	StartDt      time.Time `json:"start_dt" db:"start_dt"`
	EndDt        time.Time `json:"end_dt" db:"end_dt"`
	UpdateDt     time.Time `json:"update_dt" db:"update_dt"`
}

// Function the handler will use to create an event
func CreateEvent (db *sqlx.DB, ae ApiEvent) error {
	var err error

	// Look at this honkin' query -- we use a sub-query to grab the user_id for a given email address.
	query_str, err := db.Preparex(`INSERT INTO tinyplannr_api.event
	                                   (user_id, title, description, location, all_day, start_dt, end_dt, update_dt)
	                               VALUES ((SELECT user_id FROM tinyplannr_api.user WHERE email = $1),
	                                       $2, $3, $4, $5, $6, $7, $8)`)

	if err != nil {
		return err
	}

	query_str.QueryRowx(ae.Email, ae.Title, ae.Description, ae.Location, ae.AllDay, ae.StartDt, ae.EndDt, ae.UpdateDt)

	if err != nil {
		return err
	}

	return err
}
