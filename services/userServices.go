package services

import (
	"encoding/json"
	"errors"
	"github.com/arfan21/getprint-order/models"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var URL = os.Getenv("SERVICE_USER")

func GetUser(id uint) (map[string]interface{}, error) {
	res, err := http.Get(URL + "user/" + strconv.FormatUint(uint64(id), 10))

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, models.ErrInternalServerError
		}
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	decodeJSON := make(map[string]interface{})

	err = json.Unmarshal(body, &decodeJSON)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, errors.New("user not found")
	}

	decodeJSON["status_code"] = res.StatusCode

	return decodeJSON, nil
}
