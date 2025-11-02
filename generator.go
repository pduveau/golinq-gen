package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

type TmplAttribute struct {
	Class string
	Field string
	ATag  string
}

type TmplChild struct {
	Class    string
	Field    string
	FieldSfx string
	Space    string
	Local    string
	CTag     string
	Type     string
	TypeNew  string
}

type Field struct {
	Sep, Name string
}

type TmplClass struct {
	classType     classType
	Github        string
	Package       string
	Imports       []string
	Class         string
	ClassNs       string
	ClassTag      string
	Attributes    []TmplAttribute
	HasProperties bool
	PropertyClass string
	PropertyTag   string
	Empty         []TmplChild
	Val           []TmplChild
	Text          []TmplChild
	Children      []TmplChild
	Fields        []Field
}

func (c TmplClass) HasAttributes() bool {
	return len(c.Attributes) > 0
}

func (c TmplClass) HasChildren() bool {
	return len(c.Fields) > 0
}

func (c TmplClass) Shared() bool {
	return c.classType == Shared
}

func (c TmplClass) HasText() bool {
	return c.classType == HasText
}

func (c TmplClass) HasTextOnly() bool {
	return len(c.Fields) == 0 && c.classType == HasText
}

type TmplSwitch struct {
	Package    string
	Imports    string
	Cases      string
	Selector   string
	Github     string
	Initialize bool
}

func (c *class) createClassFile(dir string, outofns string) (err error) {
	pkg := c.Element.Gopackage

	if c.Element.ELocal == "" || c.classType > CommonPkg {
		return
	}

	var data = TmplClass{
		Package:   pkg,
		Class:     c.Element.Goname,
		ClassNs:   c.Element.Goname + "_" + c.Element.ESpace,
		Github:    github,
		ClassTag:  c.Element.ETag,
		classType: c.classType,
	}

	template := "class"

	if c.isProperties {
		template = "property"
	}

	if c.childGeneric {
		template = "generic"
	}

	for _, v := range c.Attributes {
		data.Attributes = append(data.Attributes,
			TmplAttribute{
				Class: c.Element.Goname,
				Field: v.Goname,
				ATag:  v.ETag,
			})
	}

	imports := make(map[string]bool)
	i := 0
	for _, v := range c.Children {
		sep := ""
		if i%6 == 5 {
			sep = "\n"
		}
		if v.Element.ETag == c.Element.ETag+"Pr" && template != "property" {
			data.HasProperties = true
			data.PropertyClass = v.Element.Goname
			data.PropertyTag = v.Element.ETag
		} else {
			child := TmplChild{
				Class:   c.Element.Goname,
				Field:   v.Element.Goname,
				Type:    v.Element.Goname,
				CTag:    v.Element.ETag,
				Space:   v.Element.ESpace,
				Local:   v.Element.ELocal,
				TypeNew: "New" + v.Element.Goname,
			}
			switch v.classType {
			case Text:
				child.Type = "Text"
				data.Text = append(data.Text, child)
			case RawText:
				child.Type = "RawText"
				data.Text = append(data.Text, child)
			case Empty:
				data.Empty = append(data.Empty, child)
			case Valclass:
				data.Val = append(data.Val, child)
			default:
				if v.Element.Suffix != "" {
					child.FieldSfx = v.Element.GonameShort + v.Element.Suffix
				}
				if v.Element.Gopackage != "" && v.Element.Gopackage != c.Element.Gopackage || outofns != "" {
					child.Type = v.Element.Gopackage + "." + child.Type
					child.TypeNew = v.Element.Gopackage + "." + child.TypeNew
					imports[v.Element.Gopackage] = true
				}
				data.Children = append(data.Children, child)
			}
			data.Fields = append(data.Fields, Field{Sep: sep, Name: child.Field})
			i++
		}
	}

	for v := range imports {
		data.Imports = append(data.Imports, github+"/linq/"+v)
	}

	if outofns != "" {
		data.Package = outofns
	}

	file := c.Element.ESpace + "_" + c.Element.ELocal + ".go"

	return parseFile(dir, file, template, data)
}

func (c *class) createDerivedFile(dir string) (err error) {
	file := c.Element.ESpace + "_" + c.Element.ELocal + ".go"
	return parseFile(dir, file, "aacircular", struct{ Class, Github string }{Class: c.Element.Goname, Github: github})
}

func parseInitLinq(dir string) error {
	return parseFile(dir, "initLinq.go", "initlinq", github)
}

func parseFile(dir, file, template string, data any) (err error) {
	var f *os.File

	os.MkdirAll(dir, 0775)
	path := filepath.Join(dir, file)

	f, err = os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	var buf = &bytes.Buffer{}
	var b []byte
	err = tmpl.ExecuteTemplate(buf, template, data)
	if err != nil {
		return
	}

	b, err = format.Source(buf.Bytes())
	if err != nil {
		f.Write(buf.Bytes())
		fmt.Printf("%s : %v\n", f.Name(), err)
		return nil
	}

	_, err = f.Write(b)
	return
}
