package models

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
)

type Website struct {
	Name        string  `dynamodbav:"name"`
	Url         string  `dynamodbav:"url"`
	Hash        string  `dynamodbav:"hash"`
	Subscribers []int64 `dynamodbav:"subscribers"`
}

func (w *Website) CheckChanged() bool {
	// get websites hash
	data := getWebsiteData(w.Url)
	hash := getWebsiteHash(data)

	// if hash is different, update hash and return true
	if hash != w.Hash {
		w.Hash = hash
		return true
	}

	return false
}

func getWebsiteData(url string) string {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic("Error creating request. Here's why: " + err.Error())
	}
	res, err := client.Do(req)
	if err != nil {
		panic("Error sending request. Here's why: " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error closing body. Here's why: " + err.Error())
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("Error reading body. Here's why: " + err.Error())
	}
	return string(body)
}

func getWebsiteHash(data string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		panic("Error hashing website. Here's why: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
