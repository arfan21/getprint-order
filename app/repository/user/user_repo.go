package user

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (map[string]interface{}, error)
}

type userRepository struct {
	url string
}

func NewUserRepository() UserRepository {
	url := os.Getenv("SERVICE_USER")
	return &userRepository{url}
}

func (repo *userRepository) GetUserByID(ctx context.Context, id string) (map[string]interface{}, error) {
	client := new(http.Client)

	req, err := http.NewRequestWithContext(ctx, "GET", repo.url+"/user/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	decodedJSON := make(map[string]interface{})

	err = json.Unmarshal(body, &decodedJSON)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, errors.New(decodedJSON["message"].(string))
	}

	return decodedJSON["data"].(map[string]interface{}), nil
}
