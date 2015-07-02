package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cyarie/tinyplannr-api-v2/api/settings"
)

type AppHandler struct {
	*settings.AppContext
	AuthRoute bool
	RouteName string
	H         func(*settings.AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn.H(fn.AppContext, w, r)

	log.Printf(
		"%s\t%s\t%d\t%s",
		r.Method,
		r.RequestURI,
		fn.AppContext.HandlerResp,
		fn.RouteName,
	)

	// Error-handling block. Take in a
	if err != nil {
		log.Println(err)
		switch status {
		case http.StatusNotFound:
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(settings.JsonResp{status, "Object not found. Please try again."})
		case http.StatusInternalServerError:
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(settings.JsonResp{status, "Encountered an internal server error. Please try again."})
		case http.StatusUnauthorized:
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(settings.JsonResp{status, "Login failed. Please provide your email address and password and try again."})
		case http.StatusConflict:
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(settings.JsonResp{status, "User already exists. Please try again."})
		case 422:
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(settings.JsonResp{status, "Request included malformed JSON. Please try again."})
		}
	}
}
