package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/tars47/go-sitemap-builder/sitemap"
)

var site string

func init() {
	flag.StringVar(&site, "site", "", "an url that you want to build a sitemap for (example: https://google.com)")
}

func main() {
	t1 := time.Now()
	flag.Parse()

	if site == "" || !IsUrl(site) {
		flag.Usage()
		log.Fatalf("site argument is required and should be a url")
	}

	links := sitemap.BuildXml(site)
	fmt.Printf("Time took %v\n", time.Since(t1))

	for u := range links {
		fmt.Println(u)
	}
	fmt.Println("\n", len(links))

}

func IsUrl(str string) bool {
	u, err := url.Parse(str)

	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	return true
}
