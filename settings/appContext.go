package settings

import (
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AppContext struct {
	Db					*sqlx.DB
	CookieMachine		*securecookie.SecureCookie
	HandlerResp			int
}

type JsonErr struct {
	Code		int		`json:"code"`
	Text		string	`json:"error"`
}