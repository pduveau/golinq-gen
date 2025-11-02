package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Toc struct {
	Items []struct {
		Items []struct {
			Name string `json:"name"`
			Href string `json:"href"`
		} `json:"items"`
	} `json:"items"`
}

func getDataReader(url string) io.ReadCloser {
	var f io.ReadCloser
	var e error

	f, e = os.Open("data/" + url)
	if e != nil {
		var r *http.Response
		var o *os.File

		r, e = http.Get("https://docs.dndocs.com/n/DocumentFormat.OpenXml.Linq/3.1.0/api/" + url)
		if e != nil {
			log.Fatalf("%v", e.Error())
		}
		defer r.Body.Close()

		os.MkdirAll("data", 0777)
		o, e = os.Create("data/" + url)
		if e != nil {
			log.Fatalf("%v", e.Error())
		}

		_, e = io.Copy(o, r.Body)
		if e != nil {
			log.Fatalf("%v", e.Error())
		}
		o.Close()

		f, _ = os.Open("data/" + url)
	}

	return f
}

func getToc() []string {
	var f io.ReadCloser
	var e error
	var toc Toc

	f = getDataReader("toc.json")

	dec := json.NewDecoder(f)
	e = dec.Decode(&toc)

	if e != nil {
		log.Fatalf("%v\n", e.Error())
	}

	urls := []string{}

	for _, v := range toc.Items[0].Items {
		urls = append(urls, v.Href)
	}

	return urls
}

var XmlElementRE = regexp.MustCompile(`Represents the ([:A-Za-z0-9-]+) XML (element|elements|attribute|attributes|element and attribute|elements and attribute|element and attributes|elements and attributes)\.`)
var NameSpaceRE = regexp.MustCompile(`Defines the XML namespace associated with the ([:A-Za-z0-9-]+) prefix\.`)

//var XmlElementREPublic = regexp.MustCompile(`public static readonly (XName|XNamespace) ([:A-Za-z0-9-]+)`)

func attributeToGoName(name string) string {
	if len(name) > 1 {
		return strings.ToUpper(name[0:1]) + name[1:]
	}
	return strings.ToUpper(name)
}

func parseNamespace(f io.ReadCloser, url string) {
	var doc *goquery.Document
	var e error

	elements := make(map[string]*class)
	attributes := make(map[string]*element)
	ns := ""
	doc, e = goquery.NewDocumentFromReader(f)
	if e != nil {
		log.Fatalf("%v", e)
	}
	f.Close()

	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		remarks := false
		xmlName := true
		//isNotNS := true
		type_ := ""

		c := s.Nodes[0].NextSibling
		obj := strings.Trim(s.Nodes[0].FirstChild.Data, " \t\n_")

		var cl = newClass()

		for ; c != nil && (c.Data != "h3" || c.Type != html.ElementNode); c = c.NextSibling {
			if c.Type == html.ElementNode {
				switch c.Data {
				case "h4":
					remarks = true
				case "div":
					if xmlName {
						fchild := c.FirstChild
						if fchild != nil && fchild.FirstChild != nil {
							b := fchild.FirstChild.Data
							if XmlElementRE.MatchString(b) {
								t := XmlElementRE.FindStringSubmatch(b)
								cl.Element.ETag = t[1]
								type_ = t[2]
								u := strings.Split(t[1], ":")
								cl.Element.Goname = attributeToGoName(obj)
								cl.Element.GonameShort = attributeToGoName(obj)
								cl.Element.ELocal = u[len(u)-1]
							}
							if NameSpaceRE.MatchString(b) {
								ns = NameSpaceRE.FindStringSubmatch(b)[1]
							}
							xmlName = false
						}
					}
					if remarks {
						switch type_ {
						case "attribute", "attributes":
							attributes[obj] = &parseUL(c.FirstChild, cl).Element
						case "element", "elements":
							elements[obj] = parseUL(c.FirstChild, cl)
						case "element and attribute", "elements and attribute", "element and attributes", "elements and attributes":
							elements[obj], attributes[obj] = parseDouble(c.FirstChild, cl)
						}
					}
				}
			}
		}
	})

	url2namespace[url] = ns

	xmlElements[ns] = elements
	xmlAttributes[ns] = attributes
	xmlAliases[ns] = map[string]*class{}
}

func parseDouble(in *html.Node, cl *class) (cls *class, att *element) {
	attr := cl.Element
	att = &attr
	cls = cl

	for c := in; c != nil; c = c.NextSibling {
		if c.Data == "p" {
			switch c.FirstChild.Data {
			case "As an XML element, it:":
				parseUL(c.NextSibling, cl)
			}
		}
	}

	return
}

func parseUL(in *html.Node, cl *class) *class {
	for c := in; c != nil; c = c.NextSibling {
		if c.Data == "ul" {
			// go through li
			for d := c.FirstChild; d != nil; d = d.NextSibling {
				if d.Data == "li" {
					switch d.FirstChild.Data {
					case "has the following child XML elements: ":
						parseChild(d.FirstChild, cl)
					case "has the following XML attributes: ":
						parseAttrs(d.FirstChild, cl)
					}
				}
			}
			return cl
		}
	}
	return cl
}

func parseChild(in *html.Node, cl *class) {
	for c := in; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			for _, a := range c.Attr {
				if a.Key == "href" {
					cl.Xmlchildren = append(cl.Xmlchildren, xmldata{
						Name: c.FirstChild.Data,
						Url:  strings.Split(a.Val, "#")[0],
					})
				}
			}
		}
	}
}

func parseAttrs(in *html.Node, cl *class) {
	for c := in; c != nil; c = c.NextSibling {
		if c.Data == "a" {
			for _, a := range c.Attr {
				if a.Key == "href" {
					cl.Xmlattribs = append(cl.Xmlattribs, xmldata{
						Name: c.FirstChild.Data,
						Url:  strings.Split(a.Val, "#")[0],
					})
				}
			}
		}
	}
}
