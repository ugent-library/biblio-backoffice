package render

import (
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}
	return false
}

func BadRequest(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return true
	}
	return false
}
