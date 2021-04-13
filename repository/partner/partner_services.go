package partner

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type PartnerRepository interface {
	GetPartnerByID(id uint) (map[string]interface{}, error)
}
type partnerRepository struct {
	ctx context.Context
	url string
}

func NewPartnerRepository(ctx context.Context) PartnerRepository {
	url := os.Getenv("SERVICE_PARTNER")
	return &partnerRepository{ctx, url}
}

func (repo *partnerRepository) GetPartnerByID(id uint) (map[string]interface{}, error) {

	client := new(http.Client)
	req, err := http.NewRequestWithContext(repo.ctx, "GET", repo.url+"/partner/"+strconv.FormatUint(uint64(id), 10), nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

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