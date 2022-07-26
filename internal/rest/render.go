package rest

import (
	"encoding/json"
	"net/http"
)

func RenderResponse(w http.ResponseWriter, resp interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(resp)
	if err == nil {
		_, errw := w.Write(b)
		if errw == nil {
			w.WriteHeader(status)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
}
