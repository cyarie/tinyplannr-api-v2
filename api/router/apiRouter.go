package router

import (
	"github.com/cyarie/tinyplannr-api-v2/api/handlers"
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	"github.com/gorilla/mux"
)

func ApiRouter(ac *settings.AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Index Route
	IndexFunc := handlers.AppHandler{
		ac,
		false,
		"Index",
		handlers.IndexHandler,
	}
	router.Methods("GET").Path("/").Name(IndexFunc.RouteName).Handler(&IndexFunc)

	// User Creation Route
	UserCreateFunc := handlers.AppHandler{
		ac,
		false,
		"CreateUser",
		handlers.UserCreateHandler,
	}
	router.Methods("POST").Path("/user/create").Name(UserCreateFunc.RouteName).Handler(&UserCreateFunc)

	// User Deletion Route
	UserDeleteFunc := handlers.AppHandler{
		ac,
		false,
		"DeleteUser",
		handlers.UserDeleteHandler,
	}
	router.Methods("DELETE").Path("/user/delete").Name(UserCreateFunc.RouteName).Handler(&UserDeleteFunc)

	// User Index Route
	UserIndexFunc := handlers.AppHandler{
		ac,
		false,
		"UserIndex",
		handlers.UserIndexHandler,
	}
	router.Methods("GET").Path("/user/{userId}").Name(UserIndexFunc.RouteName).Handler(&UserIndexFunc)

	// Event Creation Route
	EventCreateFunc := handlers.AppHandler{
		ac,
		false,
		"CreateEvent",
		handlers.CreateEventHandler,
	}
	router.Methods("POST").Path("/event/create").Name(UserIndexFunc.RouteName).Handler(&EventCreateFunc)

	return router
}
