package crunch

import (
	"fmt"
	"io"
	"strings"
)

type TsvWriter struct {
	parts  []string
	writer io.Writer
	index  int
}

func NewTsvWriter(w io.Writer, row *Row) *TsvWriter {
	return &TsvWriter{writer: w, parts: make([]string, row.Size())}
}

func (self *TsvWriter) Field(val string) {
	self.parts[self.index] = val
	self.index++
}

func (self *TsvWriter) End() {
	if len(self.parts) != self.index {
		panic(fmt.Sprintf("did not fill up a row: %i", len(self.parts)))
	}
	io.WriteString(self.writer, strings.Join(self.parts, "\t")+"\n")
	self.index = 0
}
