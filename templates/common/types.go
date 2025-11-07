// Generated using common template
package common

import (
	"encoding/xml"
	"reflect"
	"slices"
	"strconv"
)

type Container interface {
	MarshalXML(e *xml.Encoder, start xml.StartElement) error
	UnmarshalXML(d *xml.Decoder, start xml.StartElement) error
	SetParent(Container)
	GetParent() Container
}

type FuncXMLNameToContainer func(xml.Name, Container) Container

func XmlName(n xml.Name) string {
	if n.Space != "" {
		return n.Space + ":" + n.Local
	}
	return n.Local
}

func MarshalContainer(i Container, e *xml.Encoder) error {
	if i != nil {
		switch reflect.TypeOf(i).Kind() {
		case reflect.Ptr:
			if !reflect.ValueOf(i).IsNil() {
				return i.MarshalXML(e, xml.StartElement{})
			}
		}
	}
	return nil
}

func SetString(val ...any) *string {
	if len(val) > 0 {
		switch v := val[0].(type) {
		case string:
			return &v
		case int:
			str := strconv.Itoa(v)
			return &str
		case uint:
			str := strconv.FormatUint(uint64(v), 10)
			return &str
		case int64:
			str := strconv.FormatInt(v, 10)
			return &str
		case uint64:
			str := strconv.FormatUint(v, 10)
			return &str
		case float32:
			str := strconv.FormatFloat(float64(v), 'f', 2, 32)
			return &str
		case float64:
			str := strconv.FormatFloat(v, 'f', 2, 64)
			return &str
		}
	}
	return nil
}

func GetTyped[T any](o []Container) ([]T, []int) {
	ret := make([]T, 0)
	pos := make([]int, 0)
	for i, v := range o {
		a, ok := (v).(T)
		if ok {
			ret = append(ret, a)
			pos = append(pos, i)
		}
	}
	return ret, pos
}

func SetChild[T any](child **T, val ...any) *T {
	if len(val) > 0 {
		switch v := val[0].(type) {
		case *T:
			var t = *v
			*child = &t
			return *child
		}
	}
	var t T
	*child = &t
	return *child
}

func SetChildren[T any](o []Container, v ...Container) []Container {
	_, t_i := GetTyped[T](o)
	for _, i := range t_i {
		o = slices.Delete(o, i, i+1)
	}
	return append(o, v...)
}
