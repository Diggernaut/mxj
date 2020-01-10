package mxj

import (
	"bytes"
	"encoding/xml"
	"reflect"
)

const (
	DefaultElementTag = "element"
)

// Encode arbitrary value as XML.
//
// Note: unmarshaling the resultant
// XML may not return the original value, since tag labels may have been injected
// to create the XML representation of the value.
/*
 Encode an arbitrary JSON object.
	package main

	import (
		"encoding/json"
		"fmt"
		"github.com/clbanning/mxj"
	)

	func main() {
		jsondata := []byte(`[
			{ "somekey":"somevalue" },
			"string",
			3.14159265,
			true
		]`)
		var i interface{}
		err := json.Unmarshal(jsondata, &i)
		if err != nil {
			// do something
		}
		x, err := mxj.AnyXmlIndent(i, "", "  ", "mydoc")
		if err != nil {
			// do something else
		}
		fmt.Println(string(x))
	}

	output:
		<mydoc>
		  <somekey>somevalue</somekey>
		  <element>string</element>
		  <element>3.14159265</element>
		  <element>true</element>
		</mydoc>
*/
// Alternative values for DefaultRootTag and DefaultElementTag can be set as:
// AnyXml( v, myRootTag, myElementTag).
func AnyXml(v interface{}, tags ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.Marshal(v)
	}

	var err error
	s := new(string)
	p := new(pretty)

	var rt, et string
	if len(tags) == 1 || len(tags) == 2 {
		rt = tags[0]
	} else {
		rt = DefaultRootTag
	}
	if len(tags) == 2 {
		et = tags[1]
	} else {
		et = DefaultElementTag
	}

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">"
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndent(false, s, tag, val, p)
					}
				} else {
					err = mapToXmlIndent(false, s, et, vv, p)
				}
			default:
				err = mapToXmlIndent(false, s, et, vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		m := Map(v.(map[string]interface{}))
		b, err = m.Xml(rt)
	default:
		err = mapToXmlIndent(false, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

// Encode an arbitrary value as a pretty XML string.
// Alternative values for DefaultRootTag and DefaultElementTag can be set as:
// AnyXmlIndent( v, "", "  ", myRootTag, myElementTag).
func AnyXmlIndent(v interface{}, prefix, indent string, tags ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
	}

	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var rt, et string
	if len(tags) == 1 || len(tags) == 2 {
		rt = tags[0]
	} else {
		rt = DefaultRootTag
	}
	if len(tags) == 2 {
		et = tags[1]
	} else {
		et = DefaultElementTag
	}
	rt = checkKey(rt)
	et = checkKey(et)
	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">\n"
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndent(true, s, checkKey(tag), val, p)
					}
				} else {
					p.start = 1 // we 1 tag in
					err = mapToXmlIndent(true, s, et, vv, p)
					*s += "\n"
				}
			default:
				p.start = 0 // in case trailing p.start = 1
				err = mapToXmlIndent(true, s, et, vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		m := Map(v.(map[string]interface{}))
		b, err = m.XmlIndent(prefix, indent, rt)
	default:
		err = mapToXmlIndent(true, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}
func AnyXmlIndentByte(v interface{}, prefix, indent string, tags ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
	}
	var err error
	var buffer bytes.Buffer
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var rt, et string
	if len(tags) == 1 || len(tags) == 2 {
		rt = tags[0]
	} else {
		rt = DefaultRootTag
	}
	if len(tags) == 2 {
		et = tags[1]
	} else {
		et = DefaultElementTag
	}
	rt = checkKey(rt)
	et = checkKey(et)
	var b []byte
	switch v.(type) {
	case []interface{}:
		_, err = buffer.Write([]byte("<" + checkKey(rt) + ">\n"))
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndentByte(true, &buffer, checkKey(tag), val, p)
					}
				} else {
					p.start = 1 // we 1 tag in
					err = mapToXmlIndentByte(true, &buffer, et, vv, p)
					_, err = buffer.Write([]byte("\n"))
				}
			default:
				p.start = 0 // in case trailing p.start = 1
				err = mapToXmlIndentByte(true, &buffer, et, vv, p)
			}
			if err != nil {
				break
			}
		}
		buffer.Write([]byte("</" + rt + ">"))
		b = buffer.Bytes()
	case map[string]interface{}:
		m := Map(v.(map[string]interface{}))
		b, err = m.XmlIndentByte(prefix, indent, rt)
	default:
		err = mapToXmlIndentByte(true, &buffer, rt, v, p)
		b = buffer.Bytes()
	}
	return b, err
}
func AnyXmlIndentByteSpecial(v interface{}, prefix, indent string, tags ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
	}
	var err error
	var buffer bytes.Buffer
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var rt, et string
	if len(tags) == 1 || len(tags) == 2 {
		rt = tags[0]
	} else {
		rt = DefaultRootTag
	}
	if len(tags) == 2 {
		et = tags[1]
	} else {
		et = DefaultElementTag
	}
	rt = checkKey(rt)
	et = checkKey(et)
	var b []byte
	switch v.(type) {
	case []interface{}:
		_, err = buffer.Write([]byte("<" + checkKey(rt) + ">\n"))
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndentByteSpecial(true, &buffer, checkKey(tag), val, p, true)
					}
				} else {
					p.start = 1 // we 1 tag in
					err = mapToXmlIndentByteSpecial(true, &buffer, et, vv, p, true)
					_, err = buffer.Write([]byte("\n"))
				}
			default:
				p.start = 0 // in case trailing p.start = 1
				err = mapToXmlIndentByteSpecial(true, &buffer, et, vv, p, true)
			}
			if err != nil {
				break
			}
		}
		buffer.Write([]byte("</" + rt + ">"))
		b = buffer.Bytes()
	case map[string]interface{}:
		m := Map(v.(map[string]interface{}))
		b, err = m.XmlIndentByteSpecial(prefix, indent, rt)
	default:
		err = mapToXmlIndentByteSpecial(true, &buffer, rt, v, p, true)
		b = buffer.Bytes()
	}
	return b, err
}
