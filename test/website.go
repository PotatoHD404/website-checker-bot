package main

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
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	body := doc.Find("*").First()
	removeJunk(body)
	if w.Xpath != "" {
		body = getXpathData(body, w.Xpath)
	}
	data = body.Text()
	// write to file
	//err = ioutil.WriteFile("test.html", []byte(data), 0644)
	hash, err := getWebsiteHash(data)
	if err != nil {
		return false, err
	}

	// if hash is different, update hash and return true
	println(hash)
	if hash != w.Hash {
		w.Hash = hash
		return true, nil
	}

	return false, nil
}

func removeJunk(data *goquery.Selection) {
	// get attributes
	//println(data.Attr("class"))
	data = data.RemoveAttr("class")
	//println(data.Attr("class"))
	data = data.RemoveAttr("id")
	data.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("script") {
			s.Remove()
		} else if s.Is("style") {
			s.Remove()
		} else {
			removeJunk(s)
		}
	})
}
func getXpathData(body *goquery.Selection, xpath string) *goquery.Selection {
	// get only body
	// remove classes and ids, keep only text recursively
	// remove /html/body from xpath
	xpath = strings.ReplaceAll(xpath, "/html/body", "")

	return body.Find(xpath)
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
