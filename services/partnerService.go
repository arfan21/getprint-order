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
)



func GetPartner(id uint) (map[string]interface{}, error) {
	url := os.Getenv("SERVICE_PARTNER")

	res, err := http.Get(url + "partner/" + strconv.FormatUint(uint64(id), 10))

	if err != nil {
		return nil, models.ErrInternalServerError
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
		return nil, errors.New("partner not found")
	}

	decodeJSON["status_code"] = res.StatusCode

	return decodeJSON, nil
}