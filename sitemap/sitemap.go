package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/tars47/go-sitemap-builder/link"
)

var seen map[string]struct{}
var wg sync.WaitGroup
var mu sync.Mutex

func BuildXml(urlStr string) map[string]struct{} {
	seen = make(map[string]struct{})

	r, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Err while calling %v, err : %v\n", urlStr, err)
		return seen
	}
	defer r.Body.Close()

	reqUrl := r.Request.URL
	base := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	links := hrefs(filter(link.Parse(r.Body)), base.String())

	for _, url := range links {
		if _, ok := seen[url]; !ok {
			seen[url] = struct{}{}
			wg.Add(1)
			go get(url, base.String())
		}
	}

	wg.Wait()

	makeXml(base.Host, seen)

	return seen
}

func get(urlStr string, base string) {

	defer wg.Done()

	r, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("Err while calling %v, err : %v\n", urlStr, err)
		return
	}
	defer r.Body.Close()

	reqUrl := r.Request.URL
	cbase := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	if cbstr := cbase.String(); cbstr != base {
		fmt.Printf("Ignoring redirects to : %v\n", cbstr)
		return
	}

	links := hrefs(filter(link.Parse(r.Body)), base)

	for _, url := range links {
		mu.Lock()
		if _, ok := seen[url]; !ok {
			seen[url] = struct{}{}
			wg.Add(1)
			go get(url, base)
		}
		mu.Unlock()
	}
}

func hrefs(links map[string]struct{}, base string) []string {

	var ret []string
	for path := range links {
		switch {
		case strings.HasPrefix(path, "/"):
			ret = append(ret, base+path)
		case strings.HasPrefix(path, base):
			ret = append(ret, path)
		}
	}
	return ret
}

func filter(links []link.Link) map[string]struct{} {

	m := make(map[string]struct{})
	for _, v := range links {
		if _, ok := m[v.Href]; !ok {
			m[v.Href] = struct{}{}
		}
	}
	return m
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func makeXml(host string, seen map[string]struct{}) {
	toXml := urlset{
		Xmlns: xmlns,
	}
	for page := range seen {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	outXml := host + ".xml"

	file, err := os.Create(outXml)
	if err != nil {
		fmt.Printf("Err while creating %v, err : %v\n", outXml, err)

	}
	defer file.Close()

	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		fmt.Printf("Err while encoding to xml err : %v\n", err)

	}
}
