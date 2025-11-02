// Generated using _val template
package common

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Generic struct with val attribute
type Val struct {
	ns, tag, nsval string

	Val *string
}

func (class *Val) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var currentToken xml.Token

	class.ns = start.Name.Space
	class.tag = start.Name.Local

	for _, t := range start.Attr {
		if t.Name.Local == "val" {
			val := t.Value
			class.Val = &val
			class.nsval = t.Name.Local
			if t.Name.Space != "" {
				class.nsval = t.Name.Space + ":" + t.Name.Local
			}
		}
	}

	for {
		currentToken, err = d.Token()
		if err != nil {
			return
		}

		switch currentToken.(type) {
		case xml.StartElement:
			if err = d.Skip(); err != nil {
				return
			}
		case xml.EndElement:
			return
		}
	}
}

func (class Val) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = class.ns + ":" + class.tag

	if class.Val != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: class.nsval}, Value: *class.Val})
		return e.EncodeElement("", start)
	}
	return nil
}

func NewVal(name xml.Name) Container {
	return &Val{
		ns:    name.Space,
		tag:   name.Local,
		nsval: name.Space + ":val",
	}
}

func NewValClass(name xml.Name) *Val {
	return &Val{
		ns:    name.Space,
		tag:   name.Local,
		nsval: name.Space + ":val",
	}
}

func (class *Val) SetVal(val any) {
	switch v := val.(type) {
	case string:
		class.Val = &v
	case int:
		str := strconv.Itoa(v)
		class.Val = &str
	}
}

func NewValXml(name xml.Name, val ...any) *Val {
	if len(val) > 0 {
		ret := NewValClass(name)
		ret.SetVal(val)
		return ret
	}
	return nil
}

func NewValStr(tag string, val ...any) *Val {
	tab := strings.Split(":"+tag, ":")
	l := len(tab) - 2
	return NewValXml(xml.Name{Space: tab[l], Local: tab[l+1]}, val)
}

func (class *Val) is(name xml.Name) bool {
	return class.ns == name.Space && class.tag == name.Local
}

func GetValField(a *Val) string {
	if a != nil {
		return *a.Val
	}
	return ""
}

func GetVal(o []Container, name xml.Name) (int, string) {
	for i, v := range o {
		a, ok := (v).(*Val)
		if ok && a.is(name) {
			return i, *a.Val
		}
	}
	return -1, ""
}

func SetVal(o []Container, name xml.Name, v ...any) []Container {
	i, _ := GetVal(o, name)
	if len(v) == 0 && i > -1 {
		// remove
		l := len(o) - 1
		if l == 1 {
			return []Container{}
		}
		switch i {
		case 0:
			return o[1:]
		case l:
			return o[:l]
		}
		return append(o[:i], o[i+1:]...)
	}
	if len(v) > 0 {
		if i == -1 {
			return append(o[:i], NewValXml(name, v[0]))
		}
		o[i].(*Val).SetVal(v[0])
	}
	return o
}
