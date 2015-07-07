package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/cyarie/tinyplannr-api-v2/api/models"
	"github.com/cyarie/tinyplannr-api-v2/api/router"
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"github.com/gorilla/securecookie"
)

var (
	server *httptest.Server
)

func Setup() *settings.AppContext {
	connect_str := fmt.Sprintf("user=tinyplannr dbname=tinyplannr_test password=%s sslmode=disable", os.Getenv("TP_PW"))
	db, _ := sqlx.Connect("postgres", connect_str)

	cookie_key, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_HASH"))
	cookie_block, _ := base64.StdEncoding.DecodeString(os.Getenv("TINYPLANNR_SC_BLOCK"))

	context := &settings.AppContext{
		Db:            db,
		CookieMachine: securecookie.New(cookie_key, cookie_block),
	}

	server = httptest.NewServer(router.ApiRouter(context))
	fmt.Println("Test server running...")

	return context
}

func CreateTestUser() error {
	// This function will give us a test user to use in tests where we need one.
	var err error

	testUser := models.ApiUserCreate{
		Email:     "test@test.com",
		Password:  "faerts",
		FirstName: "Chris",
		LastName:  "Yarie",
		ZipCode:   22201,
		UpdateDt:  time.Now(),
	}

	body, err := json.Marshal(testUser)
	if err != nil {
		log.Println(err)
	}

	// Let's make the URL for the POST request, and the request itself
	req_str := fmt.Sprintf("%s/user/create", server.URL)
	req, err := http.NewRequest("POST", req_str, bytes.NewReader(body))
	http.DefaultClient.Do(req)

	return err
}

func CreateTestEvent() error {
	var err error

	testEvent := models.ApiEvent{
		Email: "test@test.com",
		Title: "Bert's Big Event",
		Description: "Gonna be a lot of fun",
		Location: "Faerts Home",
		AllDay: false,
		StartDt: time.Now(),
		EndDt: time.Now(),
		UpdateDt: time.Now(),
	}

	body, err := json.Marshal(testEvent)
	if err != nil {
		fmt.Println(err)
	}

	req_str := fmt.Sprintf("%s/event/create", server.URL)
	req, err := http.NewRequest("POST", req_str, bytes.NewReader(body))
	http.DefaultClient.Do(req)

	return err
}

func ClearDB(context *settings.AppContext) {
	// This will reset the tables and schemas in our DB to start their id sequences at 1.
	context.Db.MustExec(`ALTER SEQUENCE tinyplannr_api.user_user_id_seq RESTART;`)
	context.Db.MustExec(`ALTER SEQUENCE tinyplannr_auth.user_auth_id_seq RESTART;`)
	context.Db.MustExec(`ALTER SEQUENCE tinyplannr_api.event_event_id_seq RESTART;`)

	// Must delete in this order to keep Postgres from complaining about foreign keys.
	context.Db.MustExec(`DELETE FROM tinyplannr_api.event WHERE event_id = 1`)
	context.Db.MustExec(`DELETE FROM tinyplannr_auth.user WHERE user_id = 1`)
	context.Db.MustExec(`DELETE FROM tinyplannr_api.user WHERE email = 'test@test.com'`)

	return
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

	err := CreateTestUser()
	if err != nil {
		fmt.Println(err)
	}

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

	ClearDB(context)
}

func TestCreateUser(t *testing.T) {
	log.Println("Starting TestCreateUser...")
	context := Setup()

	defer context.Db.Close()
	defer Teardown()

	// Let's setup our post data
	testUser := models.ApiUserCreate{
		Email:     "test@test.com",
		Password:  "faerts",
		FirstName: "Chris",
		LastName:  "Yarie",
		ZipCode:   22201,
		UpdateDt:  time.Now(),
	}

	body, err := json.Marshal(testUser)
	if err != nil {
		t.Error(err)
	}

	// Let's make the URL for the POST request, and the request itself
	req_str := fmt.Sprintf("%s/user/create", server.URL)
	req, err := http.NewRequest("POST", req_str, bytes.NewReader(body))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected a 201, recieved a %d", res.StatusCode)
	}

	ClearDB(context)
}

func TestDeleteUser(t *testing.T) {
	log.Println("Starting TestDeleteUser...")
	context := Setup()

	defer context.Db.Close()
	defer Teardown()

	err := CreateTestUser()
	if err != nil {
		fmt.Println(err)
	}

	deleteUser := models.ApiUser{
		Email: "test@test.com",
	}

	delBody, err := json.Marshal(deleteUser)
	if err != nil {
		t.Error(err)
	}

	delete_string := fmt.Sprintf("%s/user/delete", server.URL)
	req, err := http.NewRequest("DELETE", delete_string, bytes.NewReader(delBody))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected a 200, received a %d", res.StatusCode)
	}

	ClearDB(context)
}

func TestCreateEvent(t *testing.T) {
	log.Println("Starting TestCreateEvent...")
	context := Setup()

	defer context.Db.Close()
	defer Teardown()

	err := CreateTestUser()
	if err != nil {
		fmt.Println(err)
	}

	// Let's create a test event
	testEvent := models.ApiEvent{
		Email: "test@test.com",
		Title: "Bert's Big Event",
		Description: "Gonna be a lot of fun",
		Location: "Faerts Home",
		AllDay: false,
		StartDt: time.Now(),
		EndDt: time.Now(),
		UpdateDt: time.Now(),
	}

	body, err := json.Marshal(testEvent)
	if err != nil {
		t.Error(err)
	}

	req_str := fmt.Sprintf("%s/event/create", server.URL)
	req, err := http.NewRequest("POST", req_str, bytes.NewReader(body))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected a 201, recieved a %d", res.StatusCode)
	}

	ClearDB(context)
}

func TestDeleteEvent
