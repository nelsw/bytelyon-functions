package model

import (
	"bytelyon-functions/pkg/util/pretty"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Document struct {
	URL  URL
	HTML HTML
	Node *html.Node
}

func NewDocument(url, body string) (*Document, error) {

	var closer io.ReadCloser
	if body != "" {
		closer = io.NopCloser(bytes.NewBufferString(body))
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var b []byte
		if b, err = io.ReadAll(closer); err != nil {
			return nil, err
		}
		body = string(b)

		closer = res.Body
	}

	node, err := html.Parse(closer)
	if err != nil {
		return nil, err
	}

	return &Document{MakeURL(url), MakeHTML(body), node}, nil
}

func (d *Document) CollectURLs() (relative, remote []URL, err error) {

	// declare the traversal function type here to enable recursion inside the ƒn
	var fn func(*html.Node)

	// define the traversal function and use comments to love on your future status
	fn = func(n *html.Node) {

		// if the given node is an anchor tag
		if n.Type == html.ElementNode && n.Data == "a" {

			// iterate tag attributes until we find the href
			for _, a := range n.Attr {

				// fail fast if not hyperlink reference
				if a.Key != "href" {
					continue
				}

				href := MakeURL(a.Val)

				if err = href.Validate(); err != nil {
					continue
				}

				// check if the anchor is valid and not a hash/mail/tel reference or email address
				if href.StartsWith("mailto:", "tel:", "#") || href.EndsWith("@"+d.URL.Domain()) {
					continue
				}

				// check if the anchor is only a path; prefix root URL if so
				if href.StartsWith("/", "?") {
					relative = append(relative, d.URL.Append(href))
					continue
				}

				// check if the anchor starts with the domain + "https://www."
				if href.Domain() == d.URL.Domain() {
					relative = append(relative, href)
					continue
				}

				remote = append(remote, href)
			}
		}

		// continue traversing every sibling per child. give em noogies.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}

	// and ... all together ... RECURSE!
	fn(d.Node)

	return
}

func (d *Document) SponsoredProducts() any {

	m := map[string]any{}

	var ƒ func(n *html.Node)

	ƒ = func(n *html.Node) {

		if n.Type == html.ElementNode && hasAttributeContains(n, "class", "top-pla-group-inner") {

			fmt.Println("found top-pla-group-inner")
			n = n.FirstChild
			for {
				domain := attributeVal(n, "data-dtld")
				url := attributeVal(n.FirstChild.NextSibling.FirstChild, "href")
				fmt.Println(domain, url)
				if n.NextSibling != nil {
					n = n.NextSibling
				} else {
					break
				}
			}

			return
		}

		// continue traversing every sibling per child. give em noogies.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ƒ(c)
		}
	}

	ƒ(d.Node)
	pretty.Println(m)
	return m
}

func (d *Document) SponsoredProduct() any {
	return nil
}

func attributeVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func hasAttributeContains(n *html.Node, key, val string) bool {
	for _, a := range n.Attr {
		if a.Key == key && strings.Contains(a.Val, val) {
			return true
		}
	}
	return false
}
