package rssfetch

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RSSTime struct {
	time.Time
}

type RSSItem struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Content string   `xml:"description"`
	PubTime RSSTime  `xml:"pubDate"`
	Link    string   `xml:"link"`
	Guid    string   `xml:"guid"`
}

type RSSFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Items   []RSSItem `xml:"channel>item"`
}

func (rt *RSSTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	// П...ц. Каждый, блин, художник: пишет как видит
	var formats []string = []string{
		"Mon, _2 Jan 2006 15:04:05 -0700",
		"Mon, _2 Jan 2006 15:04:05 MST",
	}
	d.DecodeElement(&v, &start)
	var err error = nil
	var t time.Time
	for _, f := range formats {
		t, err = time.Parse(f, v)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	*rt = RSSTime{t}
	return nil
}

func Fetch(url string) ([]RSSItem, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Bad http code: %d", resp.StatusCode))
	}
	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parse(answer)
}

func parse(s []byte) ([]RSSItem, error) {
	var feed RSSFeed
	err := xml.Unmarshal(s, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Items, nil
}
