package tinytesting

import (
	"os"
	"fmt"
	"testing"
	"net/http"
	"io/ioutil"
	"net/http/httptest"
	"encoding/base64"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gorilla/securecookie"
	"github.com/cyarie/tinyplannr-api-v2/router"
	"github.com/cyarie/tinyplannr-api-v2/settings"
	"github.com/cyarie/tinyplannr-api-v2/models"
	"encoding/json"
)

var (
	server	*httptest.Server
)

func Setup() *settings.AppContext {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr_test password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sqlx.Connect("postgres", connect_str)
	tx := db.MustBegin()

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:				db,
		Tx:				tx,
		CookieMachine:	securecookie.New(cookie_key, cookie_block),
	}

	server = httptest.NewServer(router.ApiRouter(context))
	fmt.Println("Test server running...")

	return context
}

func Teardown() {
	fmt.Println("Test server closing...")
	server.Close()
}

func TestSetupTearDown(t *testing.T) {
	log.Println("Starting TestSetupTearDown...")
	Setup()
	defer Teardown()
}

func TestIndexHandler(t *testing.T) {
	log.Println("Starting TestIndexHandler...")
	context := Setup()
	defer Teardown()
	defer context.Db.Close()

	request, err := http.NewRequest("GET", server.URL, nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected a 200, got a %d", res.StatusCode)
	}

	bs, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
	}

	if string(bs) != "WELCOME TO GORT" {
		t.Errorf("GORT WAS NOT WELCOMED: %s", string(bs))
	}


}

func TestGetUser(t *testing.T) {
	log.Println("Starting TestGetUser...")
	context := Setup()
	defer context.Db.Close()
	defer Teardown()

	// We also have to initialize the database here
	testUser := models.ApiUser{
		Email: "test@test.com",
		FirstName: "Chris",
		LastName: "Yarie",
		ZipCode: 22201,
		UpdateDt: time.Now(),
	}

	// We have to run this to reset the SERIAL sequence back to one to avoid re-building the schema each time we test
	context.Tx.MustExec(`ALTER SEQUENCE tinyplannr_api.user_user_id_seq RESTART;`)
	context.Tx.MustExec(`INSERT INTO tinyplannr_api.user (email, first_name, last_name, zip_code, update_dt)
	             VALUES ($1, $2, $3, $4, $5);`, testUser.Email, testUser.FirstName, testUser.LastName, testUser.ZipCode, testUser.UpdateDt)
	context.Tx.Commit()

	// Retrieve good ol' user number one
	req_str := fmt.Sprintf("%s/user/1", server.URL)
	fmt.Println(req_str)
	request, err := http.NewRequest("GET", req_str, nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected a 200, got a %d", res.StatusCode)
	}

	jsonResp, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Error(err)
	}

	// Let's un-marshall this JSON and test the response. I prefer this method because it lets me access individual
	// fields within the struct
	var testResp models.ApiUser

	err = json.Unmarshal([]byte(jsonResp), &testResp)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(testResp)

	// Let's clear out this test row
	context.Db.MustExec(`DELETE FROM tinyplannr_api.user WHERE email = 'test@test.com'`)
}