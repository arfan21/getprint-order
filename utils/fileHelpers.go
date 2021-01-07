package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/vincent-petithory/dataurl"
	"strings"
)

//GetFileBufferAndFileName return filename with mimetype and buffer from base64string
func GetFileBufferAndFileName(data string)([]byte, string, error){
	buffer, err := getFileBuffer(data)
	if err != nil{
		return nil, "", err
	}

	filename, err := getFileName(data)
	if err != nil{
		return nil, "", err
	}

	return buffer, filename, nil
}

//getFileBuffer decode base64string first and returning array of buffer
func getFileBuffer(data string)([]byte, error){
	dataURL, err := dataurl.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return dataURL.Data, nil
}

//getFileName create filename from unix time now and adding mimetype
func getFileName(data string)(string, error){
	mimetype, err := getMimeType(data)
	if err != nil {
		return "", err
	}

	randArrayByte := make([]byte, 5)
	if _, err := rand.Read(randArrayByte); err != nil {
		panic(err)
	}
	randString := fmt.Sprintf("%X", randArrayByte)
	return fmt.Sprintf("%s.%s",randString, mimetype), nil
}

//getMimeType first decodes the base64string and returns the mimetype of the file
func getMimeType(data string) (string, error){
	dataURL, err := dataurl.DecodeString(data)
	if err != nil {
		return "", err
	}

	return strings.Split(dataURL.MediaType.ContentType(), "/")[1], nil
}