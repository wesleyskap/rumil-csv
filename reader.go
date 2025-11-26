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
// NewReader initializes a streaming zero-copy CSV reader with applied options.
func NewReader(r io.Reader, opts ...Option) *Reader {
	cfg := defaultReaderConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	rd := &Reader{
		src:     r,
		buf:     make([]byte, cfg.BufferSize),
		scratch: make([]byte, 0, cfg.BufferSize),
		colOffs: make([]int, 0, 64),
		colLens: make([]int, 0, 64),
		cfg:     cfg,
		line:    0,
	}
	return rd
}
// Reset clears reader state and binds to a new data stream.
func (r *Reader) Reset(src io.Reader) {
	r.src = src
	r.r = 0
	r.w = 0
	r.line = 0
	r.err = nil
	r.eof = false
	r.colOffs = r.colOffs[:0]
	r.colLens = r.colLens[:0]
	r.scratch = r.scratch[:0]
}