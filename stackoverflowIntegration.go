package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Ask(id, query string) {
	inputf := "site:stackoverflow.com " + query
	res := load(fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(inputf)))
	parseResult(id, res)

}
func load(uri string) io.ReadCloser {
	client := http.Client{}
	req, ReqErr := http.NewRequest("GET", uri, nil)
	if ReqErr != nil {
		fmt.Println("Error: ", ReqErr)
	}
	req.Header.Set("User-Agent", "Nokia2700c/10.0.011 (SymbianOS/9.4; U; Series60/5.0 Opera/5.0; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/525 (KHTML, like Gecko) Safari/525 3gpp-gba")

	resp, RespErr := client.Do(req)
	CheckErr("ErrorOnBotStackOverflow: ", RespErr)
	return resp.Body
}

func parseResult(id string, body io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".b_algo a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var (
			u  string
			ok bool
		)
		if u, ok = s.Attr("href"); !ok {
			return true
		}
		if strings.Index(u, "/tagged/") != -1 {
			return true
		}
		Q := u[len("https://stackoverflow.com/questions/8927727/"):]
		Q = strings.Replace(Q, "/", "", -1)
		Q = strings.Replace(Q, "-", " ", -1)
		SendToFb(id, "Got the answer of:\n"+Q)
		parseAnswer(id, load(u))
		return false
	})

}

func parseAnswer(id string, body io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".post-text").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 1 {
			t := s.Text()
			SendToFb(id, t)
			return false
		}
		return true
	})

}
