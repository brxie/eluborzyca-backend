package handler

import (
	"net/http"
)

func TestSession(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}
