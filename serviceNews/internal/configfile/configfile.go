package configfile

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigData struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func Load(filename string) (*ConfigData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return parse(data)

}

func parse(data []byte) (*ConfigData, error) {
	rv := new(ConfigData)
	err := json.Unmarshal(data, rv)
	if err != nil {
		return nil, err
	}

	return rv, nil
}
