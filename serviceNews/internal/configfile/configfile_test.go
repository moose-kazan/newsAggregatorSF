package configfile

import (
	"os"
	"testing"
)

func Test_parse(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		cfgjson := `{
			"rss":[
			   "https://habr.com/ru/rss/hub/go/all/?fl=ru",
			   "https://habr.com/ru/rss/best/daily/?fl=ru",
			   "https://cprss.s3.amazonaws.com/golangweekly.com.xml"
			],
			"request_period": 5
		}`
		cfgdata, err := parse([]byte(cfgjson))
		if err != nil {
			t.Errorf("Can't parse config: %v", err)
			return
		}
		if cfgdata == nil {
			t.Error("Expected: *ConfigData, but nil found!")
			return
		}
		if cfgdata.RequestPeriod != 5 {
			t.Errorf("RequestPeriod: expected 5, but %v found!", cfgdata.RequestPeriod)
		}
		if len(cfgdata.Rss) == 0 {
			t.Error("No feeds found!")
		}
	})
	t.Run("Incorrect", func(t *testing.T) {
		cfgjson := "some incorect data"
		cfgdata, err := parse([]byte(cfgjson))
		if err == nil {
			t.Error("Incorrect JSON, but no error!")
		}
		if cfgdata != nil {
			t.Errorf("Expected nil, but %v found!", cfgdata)
		}
	})

}

func TestLoad(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		cfgjson := `{
			"rss":[
			   "https://habr.com/ru/rss/hub/go/all/?fl=ru",
			   "https://habr.com/ru/rss/best/daily/?fl=ru",
			   "https://cprss.s3.amazonaws.com/golangweekly.com.xml"
			],
			"request_period": 5
		}`
		f, err := os.CreateTemp("", "multirss-*")
		if err != nil {
			t.Errorf("Can't create temporary file: %v!", err)
			return
		}
		var fileName = f.Name()
		_, err = f.Write([]byte(cfgjson))
		if err != nil {
			t.Errorf("Can't write to temporary file: %v!", err)
			return
		}
		err = f.Close()
		if err != nil {
			t.Errorf("Can't close temporary file: %v!", err)
			return
		}

		cfgdata, err := Load(fileName)
		if err != nil {
			t.Errorf("Can't parse config: %v", err)
			return
		}
		if cfgdata == nil {
			t.Error("Expected: *ConfigData, but nil found!")
			return
		}
		if cfgdata.RequestPeriod != 5 {
			t.Errorf("RequestPeriod: expected 5, but %v found!", cfgdata.RequestPeriod)
		}
		if len(cfgdata.Rss) == 0 {
			t.Error("No feeds found!")
		}
		err = os.Remove(fileName)
		if err != nil {
			t.Errorf("Can't delete temporary file: %v!", err)
		}
	})
	t.Run("Incorrect", func(t *testing.T) {
		cfgjson := "some incorect data"
		f, err := os.CreateTemp("", "multirss-*")
		if err != nil {
			t.Errorf("Can't create temporary file: %v!", err)
			return
		}
		var fileName = f.Name()
		_, err = f.Write([]byte(cfgjson))
		if err != nil {
			t.Errorf("Can't write to temporary file: %v!", err)
			return
		}
		err = f.Close()
		if err != nil {
			t.Errorf("Can't close temporary file: %v!", err)
			return
		}

		cfgdata, err := Load(fileName)
		if err == nil {
			t.Error("Incorrect JSON, but no error!")
		}
		if cfgdata != nil {
			t.Errorf("Expected nil, but %v found!", cfgdata)
		}
		err = os.Remove(fileName)
		if err != nil {
			t.Errorf("Can't delete temporary file: %v!", err)
		}
	})
}
