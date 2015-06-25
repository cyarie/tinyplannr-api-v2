package settings

import (
	"database/sql"

	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

type AppContext struct {
	Db					*sql.DB
	CookieMachine		*securecookie.SecureCookie
	HandlerResp			int
}