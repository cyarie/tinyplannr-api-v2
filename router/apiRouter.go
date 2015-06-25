package router

import (
	"github.com/cyarie/tinyplannr-api-v2/handlers"
	"github.com/cyarie/tinyplannr-api-v2/settings"
	"github.com/gorilla/mux"
)

func ApiRouter(ac *settings.AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	IndexFunc := handlers.AppHandler{
		ac,
		false,
		"Index",
		handlers.IndexHandler,
	}
	router.Methods("GET").Path("/").Name(IndexFunc.RouteName).Handler(&IndexFunc)

	return router
}