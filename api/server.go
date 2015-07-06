package api

import (
	// "net/http"
	"encoding/base64"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"

	"github.com/cyarie/tinyplannr-api-v2/api/router"
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

func Api() *mux.Router {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sqlx.Connect("postgres", connect_str)
	tx := db.MustBegin()

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:            db,
		Tx:            tx,
		CookieMachine: securecookie.New(cookie_key, cookie_block),
	}

	router := router.ApiRouter(context)

	return router
}
