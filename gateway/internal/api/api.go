package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

/*
 * TODO: Когда-нибудь в светлом будущем начать читать это из настроек!
 * Ну а пока: ключ - имя сервиса, значение - пара хост:порт где он висит
 *
 * Почему одно значение, а не массив: захотим несколько копий - попросим
 * DevOps чтобы сделали несколько копий с балансером, который сам будет
 * Обрабатывать сбои и пересылать запрос на живые ноды
 */
var apiHosts = map[string]string{
	"news":     "srvnews:10010",
	"comments": "srvcomments:10020",
	"moderate": "srvmoderate:10030",
}

type Post struct {
	Id      int
	Source  int
	Title   string
	Content string
	PubTime int64
	Link    string
	Guid    string
}

type Comment struct {
	Id          int
	IdPost      int
	Content     string
	PubTime     int64
	FlagObscene bool
}

type ModerResult struct {
	Filtered bool `json:"filtered"`
}

type API struct {
	host     string
	lastCode int
}

func New(name string) (*API, error) {
	var a API
	if apiHosts[name] == "" {
		return nil, errors.New(fmt.Sprintf("Unknown service name: %s!", name))
	}
	a.host = apiHosts[name]
	return &a, nil
}

func (a *API) LastCode() int {
	return a.lastCode
}

func (a *API) httpReq(rv interface{}, host string, path string, method string, params map[string]string, reqId string) error {
	a.lastCode = 0
	//fmt.Printf("API Request: %s %s %s %v\n", method, host, path, params)
	reqUrl := "http://" + host + path

	data := url.Values{}
	for k, v := range params {
		data.Add(k, v)
	}

	if method == "GET" {
		reqUrl += "?"
		reqUrl += data.Encode()
	}
	reqBody := strings.NewReader(data.Encode())

	//fmt.Printf("API Request: %s %s %v", method, reqUrl, reqBody)
	req, err := http.NewRequest(method, reqUrl, reqBody)
	if err != nil {
		return err
	}

	if reqId != "" {
		req.Header.Add("X-Request-Id", reqId)
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	a.lastCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Bad http code: %d", resp.StatusCode))
	}

	var resBody []byte
	resBody, err = io.ReadAll(resp.Body)
	err = json.Unmarshal(resBody, &rv)

	if err != nil {
		return err
	}

	return nil
}

func (a *API) Get(rv interface{}, path string, params map[string]string, reqId string) error {
	return a.httpReq(&rv, a.host, path, "GET", params, reqId)
}

func (a *API) Post(rv interface{}, path string, params map[string]string, reqId string) error {
	return a.httpReq(&rv, a.host, path, "POST", params, reqId)
}
