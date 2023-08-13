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
	"news":     "news:10010",
	"comments": "comments:10010",
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

type API struct {
	host string
}

func httpReq(rv interface{}, host string, path string, method string, params map[string]string) error {
	reqUrl := host + path

	data := url.Values{}
	for k, v := range params {
		data.Add(k, v)
	}

	var reqBody *strings.Reader
	if method == "GET" {
		reqUrl += "?"
		reqUrl += data.Encode()
	} else if method == "POST" {
		reqBody = strings.NewReader(data.Encode())
	}

	req, err := http.NewRequest(method, reqUrl, reqBody)
	if err != nil {
		return nil
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
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

func New(name string) (*API, error) {
	var a API
	if apiHosts[name] == "" {
		return nil, errors.New(fmt.Sprintf("Unknown service name: %s!", name))
	}
	a.host = apiHosts[name]
	return &a, nil
}

func (a *API) Get(rv interface{}, path string, params map[string]string) error {
	return httpReq(&rv, a.host, path, "GET", params)
}

func (a *API) Post(rv interface{}, path string, params map[string]string) error {
	return httpReq(&rv, a.host, path, "POST", params)
}
