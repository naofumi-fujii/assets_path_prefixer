package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	//filepath := "public/adinfo.html"
	prefix := os.Args[1]
	filepath := os.Args[2]
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	var assetsPaths []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					//fmt.Println(a.Val)
					assetsPaths = append(assetsPaths, a.Val)
					break
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					//fmt.Println(a.Val)
					assetsPaths = append(assetsPaths, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	writeToFile(content, filepath, assetsPaths, prefix)
}
func writeToFile(message []byte, filepath string, assetsPaths []string, prefix string) {
	//prefix := "foo"
	var s = string(message)
	for _, a := range assetsPaths {
		s = strings.Replace(s, a, strings.Join([]string{prefix, "/", a}, ""), -1)
	}
	err := ioutil.WriteFile(filepath, []byte(s), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
