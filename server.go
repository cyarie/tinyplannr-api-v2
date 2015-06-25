package main

import (
	// "net/http"
	"database/sql"
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
	mainDb, _ := sql.Open("postgres", connect_str)

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:				mainDb,
		CookieMachine:	securecookie.New(cookie_key, cookie_block),
	}

	context.Db.Ping()

	fmt.Println("CONNECTED TO THE DB")

	defer context.Db.Close()

	router := router.ApiRouter(context)
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println("GORT IS RUNNING")
}
