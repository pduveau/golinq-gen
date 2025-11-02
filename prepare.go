package main

import (
	"fmt"
	"strings"
)

type element struct {
	ESpace      string `json:"space,omitempty"`
	ELocal      string `json:"local,omitempty"`
	ETag        string `json:"tag,omitempty"`
	Goname      string `json:"goname,omitempty"`
	GonameShort string `json:"gonameShort,omitempty"`
	Gopackage   string `json:"gopackage,omitempty"`
	Suffix      string `json:"suffix,omitempty"`
}

type xmldata struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type classType int

const (
	Standard classType = iota
	HasText
	Shared
	CommonPkg // delimiter
	Empty
	Valclass
	Text
	RawText
)

type class struct {
	Element      element
	Aliases      []string
	classType    classType
	isProperties bool
	childGeneric bool
	Xmlchildren  []xmldata           `json:"xmlChildren,omitempty"`
	Xmlattribs   []xmldata           `json:"xmlAttributes,omitempty"`
	Children     map[string]*class   `json:"children,omitempty"`
	Attributes   map[string]*element `json:"attributes,omitempty"`
}

func newClass() *class {
	return &class{
		Children:   make(map[string]*class),
		Attributes: make(map[string]*element),
	}
}

var url2namespace = make(map[string]string)

var xmlElements = make(map[string]map[string]*class)
var xmlAttributes = make(map[string]map[string]*element)
var xmlAliases = make(map[string]map[string]*class)

var packageTranslate = map[string]string{
	"a14":       "a",
	"p14":       "p",
	"p15":       "p",
	"p223":      "p",
	"p228":      "p",
	"p188":      "p",
	"c15":       "c",
	"x14":       "x",
	"x15":       "x",
	"xne":       "x",
	"xlmsforms": "x",
	"w14":       "w",
	"m":         "w",
	"o":         "w",
	"v":         "w",
}

type properties struct {
	parent, child *class
}

// func dependencies() {
// 	type key struct{ cls, chd string }
// 	deps := make(map[key]bool)
// 	for pkg, classes := range xmlElements {
// 	loop:
// 		for _, cls := range classes {
// 			if pkg == "a" &&
// 				(cls.Element.Goname == "GraphicData" || cls.Element.Goname == "Ext") {
// 				continue loop
// 			}
// 			if np, ok := packageTranslate[pkg]; ok {
// 				pkg = np
// 			}
// 			for _, child := range cls.Children {
// 				cpkg := child.Element.Gopackage
// 				if np, ok := packageTranslate[cpkg]; ok {
// 					cpkg = np
// 				}
// 				if cpkg == pkg || cpkg == "" {
// 					continue
// 				}
// 				k0 := key{cls: pkg, chd: child.Element.Gopackage}
// 				if _, ok := deps[k0]; !ok {
// 					deps[k0] = false
// 				}
// 			}
// 		}
// 	}
// 	for k := range deps {
// 		if k.chd != "common" {
// 			fmt.Printf("%s -> %s\n", k.cls, k.chd)
// 		}
// 	}
// }

func prepare() {
	var propertyClasses = make([]properties, 0)
	for namespace, classes := range xmlElements {
		for _, cls := range classes {
			cls.Element.ESpace = namespace
			cls.Element.Gopackage = namespace
			if np, ok := packageTranslate[namespace]; ok {
				cls.Element.Gopackage = np
			}

			cls.Element.Goname += "_" + namespace
		}
	}
	for namespace, classes := range xmlElements {
		for obj, cls := range classes {
			var classType = Empty
			var elem *element
			var ecls *class
			var mcls map[string]*class

			if strings.HasSuffix(cls.Element.ELocal, "Pr") && len(cls.Element.ELocal) > 2 {
				parentTag := cls.Element.ELocal[0 : len(cls.Element.ELocal)-2]
				if parent, ok := classes[parentTag]; ok {
					// perhaps we have a parent but now check that it is really a parent
					propertyClasses = append(propertyClasses, properties{
						parent: parent,
						child:  cls,
					})
				}
			}

			for _, attr := range cls.Xmlattribs {
				ns0, ok := url2namespace[attr.Url]
				if ok {
					ns := strings.Trim(ns0, "_")
					_, ok = xmlAttributes[ns]
					if ok {
						attrib := strings.Trim(attr.Name, "_")
						elem, ok = xmlAttributes[ns][attrib]
						if ok {
							xmlElements[namespace][obj].Attributes[elem.ELocal] = elem
						} else {
							fmt.Printf("Can't find attribute %s for namespace %s\n", ns, attrib)
						}
					} else {
						fmt.Printf("Can't find attribute for namespace %s\n", ns)
					}
				} else {
					fmt.Printf("Can't find attribute namespace from url %s\n", attr.Url)
				}
				classType = Standard
			}

			if _, ok := cls.Attributes["val"]; len(cls.Attributes) == 1 && ok {
				classType = Valclass
			}

			for _, child := range cls.Xmlchildren {
				ns0, ok := url2namespace[child.Url]
				ns1 := strings.Trim(ns0, "_")
				ns := ns0

				name0 := child.Name
				name1 := strings.Trim(name0, "_")
				name := name0

				if ok {
					noError := true
					mcls, ok = xmlAliases[ns0]
					if !ok {
						mcls, ok = xmlAliases[ns1]
						ns = ns1
					}
					if ok {
						ecls, ok = mcls[name0]
						if !ok {
							ecls, ok = mcls[name1]
							name = name1
						}
						if ok {
							tmp := newClass()
							tmp.Element = ecls.Element
							//tmp.Element.Goname += "_" + ns
							tmp.Element.ETag = fmt.Sprintf("%s:%s", ns, name)
							tmp.Element.ELocal = name
							tmp.Element.ESpace = ns
							tmp.Element.Suffix = attributeToGoName(name)
							ecls = tmp
						}
					}
					if !ok {
						mcls, ok = xmlElements[ns0]
						if !ok {
							mcls, ok = xmlElements[ns1]
							ns = ns1
						}
						if ok {
							ecls, ok = mcls[name0]
							if !ok {
								ecls, ok = mcls[name1]
							}
							if !ok {
								fmt.Printf("Can't find children (%s/%s) for namespace {%s/%s}\n", name0, name1, ns0, ns1)
								noError = false
							}
						} else {
							fmt.Printf("Can't find the children namespace (%s/%s)\n", ns1, ns0)
							noError = false
						}
					}
					if ok {
						xmlElements[namespace][obj].Children[ecls.Element.ETag] = ecls
						//cls.isProperties = cls.isProperties && ns == pkg
					} else {
						if child.Url == "*" {
							cls.childGeneric = true
						} else {
							if noError {
								fmt.Printf("Can't find child namespace from url %s\n", child.Url)
							}
						}
					}
				}
				classType = Standard
			}
			if cls.classType == Standard {
				cls.classType = classType
			}
		}
	}

	for _, v := range propertyClasses {
		// check if the class candidat as property is really a child of the supposed parent
		_, v.child.isProperties = v.parent.Children[v.child.Element.ETag]
	}
}
