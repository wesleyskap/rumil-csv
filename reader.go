package rumil

import (
	"io"
)

// Reader streams and parses delimited data with zero heap allocations per record.
// Fields are ordered from largest byte size to smallest for optimal CPU cache utilization.
type Reader struct {
	src        io.Reader
	buf        []byte
	scratch    []byte
	colOffs    []int
	colLens    []int
	err        error
	line       int64
	r          int
	w          int
	cfg        ReaderConfig
	currRecord Record
	eof        bool
}