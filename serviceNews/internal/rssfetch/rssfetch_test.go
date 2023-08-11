package rssfetch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_parse(t *testing.T) {
	t.Run("CorrectFeed", func(t *testing.T) {
		s := `<rss version="2.0">
			<channel>
				<title>Test</title>
				<item>
					<title>N1</title>
					<link>https://example.com/n1</link>
					<description>Desc1</description>
					<pubDate>Tue, 4 Jul 2023 00:00:00 +0000</pubDate>
					<guid>n1</guid>
				</item>
				<item>
					<title>N2</title>
					<link>https://example.com/n2</link>
					<description>Desc1</description>
					<pubDate>Sat, 01 Jul 2023 04:39:12 GMT</pubDate>
					<guid>n2</guid>
				</item>
			</channel>
		</rss>`
		items, err := parse([]byte(s))
		if err != nil {
			t.Errorf("Can't parse feed: %v", err)
		}
		itemsCount := len(items)
		if itemsCount != 2 {
			t.Errorf("Items count expected 2, got %v", itemsCount)
		}
	})
	t.Run("IncorrectFeed", func(t *testing.T) {
		s := `bla-bla-bla`
		_, err := parse([]byte(s))
		if err == nil {
			t.Error("No errors on incorrect feed!")
		}
	})
}

func TestFetch(t *testing.T) {
	t.Run("CorrectFeed", func(t *testing.T) {
		s := `<rss version="2.0">
			<channel>
				<title>Test</title>
				<item>
					<title>N1</title>
					<link>https://example.com/n1</link>
					<description>Desc1</description>
					<pubDate>Tue, 4 Jul 2023 00:00:00 +0000</pubDate>
					<guid>n1</guid>
				</item>
				<item>
					<title>N2</title>
					<link>https://example.com/n2</link>
					<description>Desc1</description>
					<pubDate>Sat, 01 Jul 2023 04:39:12 GMT</pubDate>
					<guid>n2</guid>
				</item>
			</channel>
		</rss>`
		testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, s)
		}))
		items, err := Fetch(testSrv.URL)
		if err != nil {
			t.Errorf("Can't parse feed: %v", err)
		}
		itemsCount := len(items)
		if itemsCount != 2 {
			t.Errorf("Items count expected 2, got %v", itemsCount)
		}

	})
	t.Run("IncorrectFeed", func(t *testing.T) {
		s := `bla-bla-bla`
		testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, s)
		}))
		_, err := Fetch(testSrv.URL)
		if err == nil {
			t.Error("No errors on incorrect feed!")
		}
	})
}
