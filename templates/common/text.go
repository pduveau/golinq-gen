package common

import (
	"bytes"
	"encoding/xml"
	"slices"
	"strings"
)

// Text struct
type Text struct {
	ns, tag string

	Text  string
	Space *string
}

const (
	TextSpaceDefault  = "default"
	TextSpacePreserve = "preserve"
)

func (class Text) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = class.ns + ":" + class.tag

	if class.Space != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xml:space"}, Value: *class.Space})
	}

	if err = e.EncodeElement(class.Text, start); err != nil {
		return err
	}

	return nil
}

func (class *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var buf bytes.Buffer

	class.ns = start.Name.Space
	class.tag = start.Name.Local

	for _, attr := range start.Attr {
		if attr.Name.Local == "space" {
			class.Space = &attr.Value
			break
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.CharData:
			buf.Write([]byte(elem))
		case xml.EndElement:
			if elem == start.End() {
				class.Text = buf.String()
				return nil
			}
		}
	}
}

func (class *Text) is(name xml.Name) bool {
	return class.ns == name.Space && class.tag == name.Local
}

func NewText(name xml.Name) Container {
	return NewTextClass(name)
}

func NewTextClass(name xml.Name) *Text {
	return &Text{
		ns:  name.Space,
		tag: name.Local,
	}
}

func NewTextXml(name xml.Name, val ...string) *Text {
	if len(val) == 0 {
		t := &Text{
			ns:   name.Space,
			tag:  name.Local,
			Text: val[0],
		}
		if strings.TrimSpace(val[0]) != val[0] {
			xmlSpace := "preserve"
			t.Space = &xmlSpace
		}
		return t
	}
	return nil
}

func NewTextStr(tag string, val ...string) *Text {
	if len(val) == 0 {
		tab := strings.Split(":"+tag, ":")
		l := len(tab) - 2
		t := &Text{
			ns:   tab[l],
			tag:  tab[l+1],
			Text: val[0],
		}
		if strings.TrimSpace(val[0]) != val[0] {
			xmlSpace := "preserve"
			t.Space = &xmlSpace
		}
		return t
	}
	return nil
}

func GetText(o []Container, name xml.Name) ([]string, []int) {
	ret := make([]string, 0)
	pos := make([]int, 0)
	for i, v := range o {
		a, ok := (v).(*Text)
		if ok && a.is(name) {
			ret = append(ret, a.Text)
			pos = append(pos, i)
		}
	}
	return ret, pos
}

func SetText(o []Container, name xml.Name, val ...string) []Container {
	_, t_i := GetText(o, name)
	for _, i := range t_i {
		o = slices.Delete(o, i, 1)
	}
	for _, v := range val {
		t := NewTextClass(name)
		t.Text = v
		o = append(o, t)
	}
	return o
}

// Raw text struct: text not escaped
type RawText struct {
	ns, tag string

	XMLName xml.Name
	Text    string  `xml:",innerxml"`
	Space   *string `xml:"xml:space,attr,omitempty"`
}

func (class *RawText) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var buf bytes.Buffer

	class.ns = start.Name.Space
	class.tag = start.Name.Local

	for _, attr := range start.Attr {
		if attr.Name.Local == "space" {
			class.Space = &attr.Value
			break
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		switch elem := token.(type) {
		case xml.CharData:
			buf.Write([]byte(elem))
		case xml.EndElement:
			if elem == start.End() {
				class.Text = buf.String()
				return nil
			}
		}
	}
}

func (class *RawText) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	type _rawtext RawText
	var v = *(*_rawtext)(class)

	v.XMLName.Local = class.ns + ":" + class.tag
	v.XMLName.Space = ""

	return e.Encode(v)
}

func NewRawText(name xml.Name) Container {
	return NewRawTextClass(name)
}

func NewRawTextClass(name xml.Name) *RawText {
	return &RawText{
		ns:  name.Space,
		tag: name.Local,
	}
}

func NewRawTextXml(name xml.Name, val ...string) *RawText {
	if len(val) == 0 {
		t := &RawText{
			ns:   name.Space,
			tag:  name.Local,
			Text: val[0],
		}
		if strings.TrimSpace(val[0]) != val[0] {
			xmlSpace := "preserve"
			t.Space = &xmlSpace
		}
		return t
	}
	return nil
}

func NewRawTextStr(tag string, val ...string) *RawText {
	if len(val) == 0 {
		tab := strings.Split(":"+tag, ":")
		l := len(tab) - 2
		t := &RawText{
			ns:   tab[l],
			tag:  tab[l+1],
			Text: val[0],
		}
		if strings.TrimSpace(val[0]) != val[0] {
			xmlSpace := "preserve"
			t.Space = &xmlSpace
		}
		return t
	}
	return nil
}

func (class *RawText) is(name xml.Name) bool {
	return class.ns == name.Space && class.tag == name.Local
}

func GetRawText(o []Container, name xml.Name) ([]string, []int) {
	ret := make([]string, 0)
	pos := make([]int, 0)
	for i, v := range o {
		a, ok := (v).(*RawText)
		if ok && a.is(name) {
			ret = append(ret, a.Text)
			pos = append(pos, i)
		}
	}
	return ret, pos
}

func SetRawText(o []Container, name xml.Name, val ...string) []Container {
	_, t_i := GetRawText(o, name)
	for _, i := range t_i {
		o = slices.Delete(o, i, 1)
	}
	for _, v := range val {
		t := NewRawTextClass(name)
		t.Text = v
		o = append(o, t)
	}
	return o
}

// text fragment when interleave with children (i.e. x:t)
type TextFragment string

// Container compatibility
func (class TextFragment) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return nil
}

// Container compatibility
func (class *TextFragment) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	return nil
}

// support parameter type string and xml.CharData
func NewTextFragment(in any) *TextFragment {
	switch v := in.(type) {
	case string:
		t := TextFragment(v)
		return &t
	case xml.CharData:
		t := TextFragment(v)
		return &t
	}
	return nil
}

func GetTextFragment(o []Container) ([]string, []int) {
	g, p := GetTyped[*TextFragment](o)
	r := make([]string, 0, len(g))
	for _, v := range g {
		r = append(r, string(*v))
	}
	return r, p
}
