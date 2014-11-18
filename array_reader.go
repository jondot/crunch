package crunch

import (
	"errors"
)

type ArrayReader struct {
	vals []string
}

func NewArrayReader(vals []string) *ArrayReader {
	return &ArrayReader{vals: vals}
}

func (self *ArrayReader) QueryIndex(i int) (interface{}, error) {
	return self.vals[i], nil
}

func (self *ArrayReader) Query(parts []string) (interface{}, error) {
	return nil, errors.New("Query not supported in ArrayReader, use QueryIndex.")
}
