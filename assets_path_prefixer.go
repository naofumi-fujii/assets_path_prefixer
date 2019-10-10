package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func GetFileContentType1(filepath string) string {

	// Open File
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get the content
	contentType, err := GetFileContentType(f)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Content Type: " + contentType)
	return contentType
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func main() {
	//filepath := "public/adinfo.html"
	prefix := os.Args[1]
	filepath := os.Args[2]
	if GetFileContentType1(filepath) != "text/html; charset=utf-8" {
		log.Fatal("ContenType should be text/html; charset=utf-8")
	}

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
				if a.Key == "src" && isExternalURL(a.Val) {
					//fmt.Println(a.Val)
					assetsPaths = append(assetsPaths, a.Val)
					break
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" && isExternalURL(a.Val) {
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
func isExternalURL(assetPath string) bool {
	return (!strings.HasPrefix(assetPath, "https") || !strings.HasPrefix(assetPath, "http"))
}
func writeToFile(message []byte, filepath string, assetsPaths []string, prefix string) {
	//prefix := "foo"
	var s = string(message)
	//fmt.Println(assetsPaths)
	for _, a := range uniqArr(assetsPaths) {
		s = strings.Replace(s, a, strings.Join([]string{prefix, "/", a}, ""), -1)
	}
	err := ioutil.WriteFile(filepath, []byte(s), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func uniqArr(arr []string) []string {
	//arr := [string]{"a", "b", "c", "a"}
	m := make(map[string]bool)
	uniq := []string{}

	for _, ele := range arr {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}

	//fmt.Printf("%v", uniq) // ["a", "b", "c"]
	return uniq
}
