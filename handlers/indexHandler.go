package handlers

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/cyarie/tinyplannr-api-v2/settings"
)

func IndexHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, "WELCOME TO GORT")

	ac.HandlerResp = http.StatusOK
	return 200, nil
}