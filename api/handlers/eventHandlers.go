package handlers

import (
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"net/http"
	"github.com/cyarie/tinyplannr-api-v2/api/models"
	"io/ioutil"
	"io"
	"log"
	"encoding/json"
)

func CreateEventHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error
	var event models.ApiEvent

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		ac.HandlerResp = http.StatusInternalServerError
		return ac.HandlerResp, err
	}

	err = json.Unmarshal(body, &event)
	if err != nil {
		ac.HandlerResp = 422
		return ac.HandlerResp, err
	}

	err = models.CreateEvent(ac.Db, event)
	if err != nil {
		ac.HandlerResp = 422
		return ac.HandlerResp, err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ac.HandlerResp = http.StatusCreated
	w.WriteHeader(ac.HandlerResp)
	err = json.NewEncoder(w).Encode(settings.JsonResp{ac.HandlerResp, "New event created."})
	if err != nil {
		ac.HandlerResp = http.StatusInternalServerError
		return ac.HandlerResp, err
	}
	return ac.HandlerResp, nil
}
