package api

import (
	// "net/http"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"fmt"
	"os"
	"encoding/base64"

	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"github.com/gorilla/securecookie"
	"github.com/cyarie/tinyplannr-api-v2/api/router"
	"github.com/gorilla/mux"
)

func Api() *mux.Router {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sqlx.Connect("postgres", connect_str)
	tx := db.MustBegin()

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:				db,
		Tx:				tx,
		CookieMachine:	securecookie.New(cookie_key, cookie_block),
	}

	router := router.ApiRouter(context)

	return router
}
