package handlers

import (
	"net/http"

	"github.com/cyarie/tinyplannr-api-v2/settings"
	"fmt"
)

type AppHandler struct {
	*settings.AppContext
	AuthRoute 			bool
	RouteName 			string
	H					func(*settings.AppContext, http.ResponseWriter, *http.Request) (int, error)
}


func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := fn.H(fn.AppContext, w, r)
	if err != nil {
		fmt.Println("faerts")
	} else {
		fmt.Println("faerts")
	}

	fmt.Println(status)
}