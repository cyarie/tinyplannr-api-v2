package router

import (
	"github.com/cyarie/tinyplannr-api-v2/api/handlers"
	"github.com/cyarie/tinyplannr-api-v2/api/settings"
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

	UserCreateFunc := handlers.AppHandler{
		ac,
		false,
		"CreateUser",
		handlers.UserCreateHandler,
	}
	router.Methods("POST").Path("/user/create").Name(UserCreateFunc.RouteName).Handler(&UserCreateFunc)

	UserDeleteFunc := handlers.AppHandler{
		ac,
		false,
		"DeleteUser",
		handlers.UserDeleteHandler,
	}
	router.Methods("DELETE").Path("/user/delete").Name(UserCreateFunc.RouteName).Handler(&UserDeleteFunc)

	UserIndexFunc := handlers.AppHandler{
		ac,
		false,
		"UserIndex",
		handlers.UserIndexHandler,
	}
	router.Methods("GET").Path("/user/{userId}").Name(UserIndexFunc.RouteName).Handler(&UserIndexFunc)

	return router
}
