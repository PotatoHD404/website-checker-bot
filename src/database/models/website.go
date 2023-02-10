package models

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Website struct {
	Name        string  `dynamodbav:"name"`
	Url         string  `dynamodbav:"url"`
	Hash        string  `dynamodbav:"hash"`
	Xpath       string  `dynamodbav:"xpath"`
	Subscribers []int64 `dynamodbav:"subscribers"`
}

func NewWebsite(name string, url string, xpath string) Website {
	return Website{
		Name:        name,
		Url:         url,
		Hash:        "",
		Xpath:       xpath,
		Subscribers: make([]int64, 0),
	}
}

func (w *Website) CheckChanged() (bool, error) {
	// get websites hash
	data, err := getWebsiteData(w.Url)
	if err != nil {
		return false, err
	}
	// if xpath is not empty, get xpath
	if w.Xpath != "" {
		data, err = getXpathData(data, w.Xpath)
		if err != nil {
			return false, err
		}
	}
	hash, err := getWebsiteHash(data)
	if err != nil {
		return false, err
	}

	// if hash is different, update hash and return true
	if hash != w.Hash {
		w.Hash = hash
		return true, nil
	}

	return false, nil
}

func getXpathData(data string, xpath string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return doc.Find(xpath).Text(), nil
}

func getWebsiteData(url string) (string, error) {
	method := "GET"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error closing body. Here's why: " + err.Error())
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getWebsiteHash(data string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}
