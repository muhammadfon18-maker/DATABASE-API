package helper

import (
	"errors"
	"net/http"
)

var UserNotFound = errors.New("user is not found")

var CarNotFound = errors.New("car is not found")

var Invalid = errors.New("invalid operation")

func WriteError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
