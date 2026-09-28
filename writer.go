package main

import (
	"bufio"
	"fmt"
	"io"
	
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: bufio.NewWriter(w)}
}

func (w *Writer) Write(v *Value) {
	var reply string

	switch v.typ {
	case STRING:
		reply = fmt.Sprintf("%s%s\r\n", v.typ, v.str)
	case BULK:
		reply = fmt.Sprintf("%s%d\r\n%s\r\n", v.typ, len(v.bulk), v.bulk)
	case ERROR:
		reply = fmt.Sprintf("%s%s\r\n", v.typ, v.err)
	case NULL:
		reply = "$-1\r\n"

	}
    
	w.writer.Write([]byte(reply))
	w.writer.(*bufio.Writer).Flush()
}