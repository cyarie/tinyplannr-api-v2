package main

import (
	// "net/http"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"fmt"
	"os"
	"encoding/base64"

	"github.com/cyarie/tinyplannr-api-v2/settings"
	"github.com/gorilla/securecookie"
	"github.com/cyarie/tinyplannr-api-v2/router"
	"log"
	"net/http"
)

func main() {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr password=%s sslmode=disable", os.Getenv("TP_PW"))
	mainDb, _ := sqlx.Connect("postgres", connect_str)

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:				mainDb,
		CookieMachine:	securecookie.New(cookie_key, cookie_block),
	}

	router := router.ApiRouter(context)
	log.Println(http.ListenAndServe(":8080", router))
}
