package handlers

import (
	"fmt"
	"net/http"

	"github.com/cyarie/tinyplannr-api-v2/api/settings"
	_ "github.com/lib/pq"
)

func IndexHandler(ac *settings.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, "WELCOME TO GORT")

	ac.HandlerResp = http.StatusOK
	return 200, nil
}
