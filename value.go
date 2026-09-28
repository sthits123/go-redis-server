package main

import (
	"fmt"
	"io"
	"strconv"
)

type ValueType string

const (
	ARRAY   ValueType = "*"
	BULK    ValueType = "$"
	STRING  ValueType = "+"
	INTEGER ValueType = ":"
	ERROR   ValueType = "-"
	NULL    ValueType = ""
)

type Value struct {
	typ   ValueType
	bulk  string
	str   string
	num   int
	err   string
	array []Value
}


func (v *Value) readArray(r io.Reader) (error){
	
	buf:=make([]byte, 4)
	_,err:=	r.Read(buf)
	if err != nil {
		return err
	}
	arrLen, err := strconv.Atoi(string(buf[1]))
	if err != nil {
		return err
	}

	for range arrLen {
		bulk:= v.readBulk(r)
		v.array = append(v.array, bulk)
	}
   
	return nil
	
}

func (v *Value) readBulk(r io.Reader) (Value) {
	buf:=make([]byte, 4)
	r.Read(buf)

	n, err := strconv.Atoi(string(buf[1]))
	if err != nil {
		fmt.Println(err)
		return Value{}
	}

	bulkbuf:= make([]byte, n+2)
	r.Read(bulkbuf)
	bulk := string(bulkbuf[:n])
	return Value{typ: BULK, bulk: bulk}
}