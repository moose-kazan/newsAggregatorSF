package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_httpReq(t *testing.T) {
	type TestType struct {
		Method string
		Num    int
		Url    string
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json_line, err := json.Marshal(TestType{
			Method: r.Method,
			Num:    4,
			Url:    r.RequestURI,
		})

		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(json_line)
	}))
	defer srv.Close()

	srvHost := strings.Replace(srv.URL, "http://", "", 1)

	t.Run("GET", func(t *testing.T) {
		var val TestType
		var a API
		err := a.httpReq(&val, srvHost, "/test", "GET", map[string]string{"n": "5"}, "")
		if err != nil {
			t.Errorf("HTTP Error: %v", err)
		}
		if val.Method != "GET" {
			t.Errorf("val.Method: expected \"GET\", got: \"%v\"", val.Method)
		}
		if val.Url != "/test?n=5" {
			t.Errorf("val.Url: expected \"/test?n=5\", got: \"%v\"", val.Num)
		}
		if val.Num != 4 {
			t.Errorf("val.Num: expected \"4\", got: \"%v\"", val.Num)
		}
	})

	t.Run("POST", func(t *testing.T) {
		var val TestType
		var a API
		err := a.httpReq(&val, srvHost, "/test", "POST", map[string]string{"n": "5"}, "")
		if err != nil {
			t.Errorf("HTTP Error: %v", err)
		}
		if val.Method != "POST" {
			t.Errorf("val.Method: expected \"POST\", got: \"%v\"", val.Method)
		}
		if val.Url != "/test" {
			t.Errorf("val.Url: expected \"/test\", got: \"%v\"", val.Num)
		}
		if val.Num != 4 {
			t.Errorf("val.Num: expected \"4\", got: \"%v\"", val.Num)
		}
	})
}

func TestAPI_Get(t *testing.T) {
	type TestType struct {
		Method string
		Num    int
		Url    string
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json_line, err := json.Marshal(TestType{
			Method: r.Method,
			Num:    4,
			Url:    r.RequestURI,
		})

		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(json_line)
	}))
	defer srv.Close()

	srvHost := strings.Replace(srv.URL, "http://", "", 1)
	apiHosts["testGet"] = srvHost

	t.Run("Main", func(t *testing.T) {
		api, err := New("testGet")
		if err != nil {
			t.Errorf("Can't create API object: %v", err)
			return
		}
		var val TestType
		err = api.Get(&val, "/test", map[string]string{"n": "5"}, "")
		if err != nil {
			t.Errorf("HTTP Error: %v", err)
		}
		if val.Method != "GET" {
			t.Errorf("val.Method: expected \"GET\", got: \"%v\"", val.Method)
		}
		if val.Url != "/test?n=5" {
			t.Errorf("val.Url: expected \"/test?n=5\", got: \"%v\"", val.Num)
		}
		if val.Num != 4 {
			t.Errorf("val.Num: expected \"4\", got: \"%v\"", val.Num)
		}
	})
}

func TestAPI_Post(t *testing.T) {
	type TestType struct {
		Method string
		Num    int
		Url    string
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json_line, err := json.Marshal(TestType{
			Method: r.Method,
			Num:    4,
			Url:    r.RequestURI,
		})

		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(json_line)
	}))
	defer srv.Close()

	srvHost := strings.Replace(srv.URL, "http://", "", 1)
	apiHosts["testPost"] = srvHost

	t.Run("Main", func(t *testing.T) {
		api, err := New("testPost")
		if err != nil {
			t.Errorf("Can't create API object: %v", err)
			return
		}
		var val TestType
		err = api.Post(&val, "/test", map[string]string{"n": "5"}, "")
		if err != nil {
			t.Errorf("HTTP Error: %v", err)
		}
		if val.Method != "POST" {
			t.Errorf("val.Method: expected \"POST\", got: \"%v\"", val.Method)
		}
		if val.Url != "/test" {
			t.Errorf("val.Url: expected \"/test\", got: \"%v\"", val.Num)
		}
		if val.Num != 4 {
			t.Errorf("val.Num: expected \"4\", got: \"%v\"", val.Num)
		}
	})
}
