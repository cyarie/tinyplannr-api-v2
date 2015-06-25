package models

import (
	"time"
)

// Data struct for ApiUser
type ApiUser struct {
	UserId		int64 		`json:"user_id" sql:"user_id"`
	Email		string		`json:"email" sql:"email"`
	FirstName	string		`json:"first_name" sql:"first_name"`
	LastName	string		`json:"last_name" sql:"last_name"`
	ZipCode		int64		`json:"zip_code" sql:"zip_code"`
	IsActive	bool		`json:"is_active" sql:"is_active"`
	CreateDt	time.Time	`json:"create_dt" sql:"create_dt"`
	UpdateDt	time.Time	`json:"update_dt" sql:"update_dt"`
}


