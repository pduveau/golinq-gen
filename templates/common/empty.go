// Generated using _empty template
package common

import (
	"encoding/xml"
	"strings"
)

// Generic Empty struct
type Empty struct {
	ns, tag string
}

func (class *Empty) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var currentToken xml.Token

	class.ns = start.Name.Space
	class.tag = start.Name.Local

	for {
		currentToken, err = d.Token()
		if err != nil {
			return
		}

		switch currentToken.(type) {
		case xml.EndElement:
			return
		default:
			if err = d.Skip(); err != nil {
				return
			}
		}
	}
}

func (class Empty) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = class.ns + ":" + class.tag

	err = e.EncodeElement("", start)

	return
}

func NewEmpty(name xml.Name) Container {
	return NewEmptyClass(name)
}

func NewEmptyClass(name xml.Name) *Empty {
	return &Empty{
		ns:  name.Space,
		tag: name.Local,
	}
}

func NewEmptyXml(name xml.Name, val ...bool) *Empty {
	if len(val) == 0 || val[0] {
		return &Empty{
			ns:  name.Space,
			tag: name.Local,
		}
	}
	return nil
}

func NewEmptyStr(tag string, val ...bool) *Empty {
	if len(val) == 0 || val[0] {
		tab := strings.Split(":"+tag, ":")
		l := len(tab) - 2
		return &Empty{
			ns:  tab[l],
			tag: tab[l+1],
		}
	}
	return nil
}

func (class *Empty) Is(name xml.Name) bool {
	return class.ns == name.Space && class.tag == name.Local
}

func GetEmpty(o []Container, name xml.Name) int {
	for i, v := range o {
		a, ok := (v).(*Empty)
		if ok && a.Is(name) {
			return i
		}
	}
	return -1
}

func SetEmpty(o []Container, name xml.Name, v ...bool) []Container {
	i := GetEmpty(o, name)
	if len(v) > 0 && !v[0] && i > -1 {
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
	if len(v) == 0 || !v[0] && i == 0 {
		return append(o[:i], NewEmpty(name))
	}
	return o
}
