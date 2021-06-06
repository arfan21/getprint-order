package media

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type MediaRepository interface {
	DeleteFile(ctx context.Context, url string) error
}

type mediaRepository struct {
	url string
}

func NewMediaRepository() MediaRepository {
	meduaServiceUrl := os.Getenv("SERVICE_MEDIA")
	return &mediaRepository{meduaServiceUrl}
}

func (repo *mediaRepository) DeleteFile(ctx context.Context, url string) error {
	client := new(http.Client)
	req, err := http.NewRequestWithContext(ctx, "DELETE", repo.url+"/partner", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	query := req.URL.Query()
	query.Set("url", url)
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	mediaRes := new(MediaResponse)

	err = json.Unmarshal(body, mediaRes)
	if err != nil {
		return err
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return errors.New(mediaRes.Message.(string))
	}

	return nil
}

