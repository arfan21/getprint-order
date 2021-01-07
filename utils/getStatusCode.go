package utils

import (
	"github.com/arfan21/getprint-order/models"
	"net/http"
	"strings"
)

func GetStatusCode(err error) int {
	if strings.Contains(err.Error(), "Duplicate") {
		return http.StatusConflict
	}
	if strings.Contains(err.Error(), "not found") {
		return http.StatusNotFound
	}

	switch err {
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
