package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/cyarie/tinyplannr-api-v2/settings"
	"github.com/cyarie/tinyplannr-api-v2/models"
)

func UserIndexHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	vars := mux.Vars(r)

	var userId int64

	if userId, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		if err := json.NewEncoder(w).Encode(settings.JsonErr{http.StatusInternalServerError, "Encountered a server error. Please try again"}); err != nil {
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
