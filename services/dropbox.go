package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Dropbox interface {
	Upload(filename string, buffer []byte) (map[string]interface{}, error)
	CreateSharedLink(path string) (string,error)
	Delete(path string) error
}

type dropbox struct {
	token string
}

func NewDropbox() Dropbox{
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	return &dropbox{token: token}
}

func (dbx dropbox) Upload(filename string, buffer []byte) (map[string]interface{}, error){
	url := "https://content.dropboxapi.com/2/files/upload"
	dropboxApiArgs := "{\"path\": \"/getprint/" +filename+"\",\"mode\": \"add\",\"autorename\": true,\"mute\": false,\"strict_conflict\": false}"
	payload := new(bytes.Buffer)
	payload.Write(buffer)

	client := new(http.Client)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil{
		return nil,err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", dbx.token))
	req.Header.Add("Dropbox-API-Arg",dropboxApiArgs)
	req.Header.Add("Content-Type","application/octet-stream")
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil,err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,err
	}

	var resJSON map[string]interface{}

	err = json.Unmarshal(resBody, &resJSON)
	if err != nil {
		return nil,err
	}

	return resJSON, nil
}

func (dbx dropbox) CreateSharedLink(path string) (string,error){
	url := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"
	payload := []byte(`{"path":"`+path+`"}`)
	client := new(http.Client)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil{
		return "",err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", dbx.token))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return "",err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "",err
	}

	var resJSON map[string]interface{}

	err = json.Unmarshal(resBody, &resJSON)
	if err != nil {
		return "",err
	}

	return resJSON["url"].(string), nil
}

func (dbx dropbox) Delete(path string) error{
	url := "https://api.dropboxapi.com/2/files/delete_v2"
	payload := []byte(`{"path":"`+path+`"}`)

	client := new(http.Client)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil{
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", dbx.token))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	return nil
}