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
	GetPartnerByID(id uint) (*PartnerResponse, error)
}
type partnerRepository struct {
	ctx context.Context
	url string
}

func NewPartnerRepository(ctx context.Context) PartnerRepository {
	url := os.Getenv("SERVICE_PARTNER")
	return &partnerRepository{ctx, url}
}

func (repo *partnerRepository) GetPartnerByID(id uint) (*PartnerResponse, error) {
	client := new(http.Client)
	req, err := http.NewRequestWithContext(repo.ctx, "GET", repo.url+"/partner/"+strconv.FormatUint(uint64(id), 10), nil)
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

	partnerRes := new(PartnerResponse)

	err = json.Unmarshal(body, partnerRes)

	if err != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, errors.New(partnerRes.Message.(string))
	}

	return partnerRes, nil
}
