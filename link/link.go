package link

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) []Link {

	doc, err := html.Parse(r)
	if err != nil {
		log.Fatalln(err)
	}

	return _parse(doc)
}

func _parse(n *html.Node) []Link {

	var nAnchor *html.Node
	var sTxt string
	var links []Link
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			nAnchor = n
		}

		if nAnchor != nil && n.Type == html.TextNode {
			sTxt += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n == nAnchor {
			link := Link{Text: strings.Join(strings.Fields(sTxt), " ")}
			for _, v := range n.Attr {
				if v.Key == "href" {
					link.Href = v.Val
				}
			}
			links = append(links, link)
			sTxt = ""
			nAnchor = nil
		}
	}
	f(n)

	return links
}
