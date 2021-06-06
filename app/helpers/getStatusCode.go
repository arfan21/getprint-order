package helpers

import (
	"net/http"
	"strings"

	models2 "github.com/arfan21/getprint-order/app/models"
)

func GetStatusCode(err error) int {
	if strings.Contains(err.Error(), "Duplicate") {
		return http.StatusConflict
	}
	if strings.Contains(err.Error(), "not found") {
		return http.StatusNotFound
	}

	switch err {
	case models2.ErrConflict:
		return http.StatusConflict
	case models2.ErrNotFound:
		return http.StatusNotFound
	case models2.ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
