package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
    "io/ioutil"
    "io"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"github.com/cyarie/tinyplannr-api-v2/api/models"
	"github.com/lib/pq"
)

// Grabs the stored data for a given user
func UserIndexHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	vars := mux.Vars(r)

	var userId int64

	if userId, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		if err := json.NewEncoder(w).Encode(settings.JsonResp{http.StatusInternalServerError, "Encountered a server error. Please try again"}); err != nil {
			log.Println(err)
			ac.HandlerResp = http.StatusInternalServerError
			return 500, err
		}
	}

	user, err := models.GetUserData(ac.Db, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			ac.HandlerResp = http.StatusNotFound
			log.Println(err)
			return 404, err
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		ac.HandlerResp = http.StatusInternalServerError
		return 500, err
	}

	ac.HandlerResp = http.StatusOK
	return 200, err


}

// Creates a user in the database, both in the API schema and the Auth schema
func UserCreateHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	var user models.ApiUserCreate

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Println(err)
		ac.HandlerResp = http.StatusInternalServerError
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		ac.HandlerResp = 422
		return ac.HandlerResp, err
	}

	err = models.CreateUser(ac.Db, user)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Name() == "unique_violation" {
			ac.HandlerResp = http.StatusConflict
			return ac.HandlerResp, err
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ac.HandlerResp = http.StatusCreated
	w.WriteHeader(ac.HandlerResp)
	err = json.NewEncoder(w).Encode(settings.JsonResp{ac.HandlerResp, "New user created."})
	if err != nil {
		ac.HandlerResp = http.StatusInternalServerError
		return ac.HandlerResp, err
	}
	return ac.HandlerResp, nil
}

// Deletes a user from the database, again, both in the API and Auth schemas
func UserDeleteHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	var UserDel struct {
		Email    string
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Printf("Found an error in the UserDeleteHandler while parsing the body: %v", err)
		ac.HandlerResp = http.StatusInternalServerError
		return ac.HandlerResp, err
	}

	err = json.Unmarshal(body, &UserDel)

	if err != nil {
		ac.HandlerResp = 422
		log.Printf("Encountered some malformed JSON: %v", err)
		return ac.HandlerResp, err
	}

	err = models.DeleteUser(ac.Db, UserDel.Email)
	if err != nil {
		log.Printf("Encountered an error from the database while deleting a user: %v", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ac.HandlerResp = 204
	w.WriteHeader(ac.HandlerResp)
	err = json.NewEncoder(w).Encode(settings.JsonResp{ac.HandlerResp, "New user created."})
	if err != nil {
		ac.HandlerResp = http.StatusInternalServerError
		log.Printf("Server error in sending back a UserDelete response: %v", err)
		return ac.HandlerResp, err
	}

	return ac.HandlerResp, nil
}
