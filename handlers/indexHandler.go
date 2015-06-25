package handlers

import (
	"fmt"
	"net/http"

	"github.com/cyarie/tinyplannr-api-v2/settings"
)

func IndexHandler(a *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, "WELCOME TO GORT")
	return 200, nil
}