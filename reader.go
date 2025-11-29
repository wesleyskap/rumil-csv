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
// fill loads incoming bytes from the underlying reader into the sliding buffer.
func (r *Reader) fill() error {
	if r.r > 0 {
		n := copy(r.buf, r.buf[r.r:r.w])
		r.w = n
		r.r = 0
	}
	if r.w >= len(r.buf) {
		r.growBuffer()
	}
	n, err := r.src.Read(r.buf[r.w:])
	r.w += n
	if err != nil {
		if err == io.EOF {
			r.eof = true
			return nil
		}
		r.err = err
		return err
	}
	return nil
}
// growBuffer expands the internal buffer up to configured maximum record size.
func (r *Reader) growBuffer() {
	newCap := len(r.buf) * 2
	if newCap > r.cfg.MaxRecordSize {
		newCap = r.cfg.MaxRecordSize
	}
	newBuf := make([]byte, newCap)
	copy(newBuf, r.buf[:r.w])
	r.buf = newBuf
}
// consumeNewline increments line counter and moves past CRLF or LF.
func (r *Reader) consumeNewline() {
	r.line++
	if r.r < r.w && r.buf[r.r] == '\r' {
		r.r++
	}
	if r.r < r.w && r.buf[r.r] == '\n' {
		r.r++
	}
}