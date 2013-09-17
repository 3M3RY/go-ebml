// Copyright © 2013 Emery Hemingway
// Released under the terms of the GNU Public License version 3

/*
Package ebml marshals and unmarshals Go structs to and from
the Exensible Binary Markup Langauge.
*/
package ebml

import (
	"fmt"
	"reflect"
	"strconv"
)

// Header is a struct for encoding and decoding EBML streams.
//
// The following could be a valid top-level struct for
// representing Matroska streams:
//	type Matroska struct {
//		ebml.Header
//		Segment     []Segment `ebml:"18538067"`
//	}
//
// You will however need to populate field values in Header
// to form a valid EBML document.
type Header struct {
	// BUG(Emery): Setting default and fallback values in
	// struct field tags is not supported.
	EbmlId             Id     `ebml:"1a45dfa3"`
	EBMLVersion        uint8  `ebml:"4286"`
	EBMLReadVersion    uint8  `ebml:"42f7"`
	EBMLMaxIDLength    uint8  `ebml:"42f2"`
	EBMLMaxSizeLength  uint8  `ebml:"42f3"`
	DocType            string `ebml:"4282"`
	DocTypeVersion     uint8  `ebml:"4287"`
	DocTypeReadVersion uint8  `ebml:"4285"`
}

// Id is a type that identifies an ebml element.
type Id uint64

// Bytes returns an Id in byte slice form
func (id Id) Bytes() []byte {
	var l int
	switch {
	case id < 0x81:
		panic(fmt.Sprintf("invalid element ID %s", id))
	case id < 0x4001:
		l = 1
	case id < 0x200001:
		l = 2
	case id < 0x10000001:
		l = 3
	case id < 0x0800000001:
		l = 4
	default:
		panic(fmt.Sprintf("invalid element ID %s", id))
	}
	b := make([]byte, l)
	for l > 1 {
		l--
		b[l] = byte(id)
		id >>= 8
	}
	b[0] = byte(id)
	return b
}

// NewIdFromString is a convience fuction that returns an Id
// for a string of hexadecimal.
func NewIdFromString(s string) (Id, error) {
	x, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return Id(x), nil
}

// String implements the Stringer interface, returning Id in
// hexadecimal form.
func (id Id) String() string {
	return strconv.FormatUint(uint64(id), 16)
}

// getId parses a struct for "EbmlId"
func getId(v reflect.Value) (id Id) {
	v = reflect.Indirect(v)
	t := v.Type()
	if f, ok := t.FieldByName("EbmlId"); ok {
		id = v.FieldByIndex(f.Index).Interface().(Id)
		if id == 0 {
			var err error
			id, err = NewIdFromString(f.Tag.Get("ebml"))
			if err != nil {
				panic("cannot resolve EBML Id for " + t.Name() + ", " + err.Error())
			}
		}
	} else {
		panic("cannot resolve EBML Id for " + t.Name())
	}
	return
}
