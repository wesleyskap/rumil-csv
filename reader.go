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
// skipCommentOrEmpty advances past comments and blank lines.
func (r *Reader) skipCommentOrEmpty() bool {
	if r.r >= r.w {
		return false
	}
	b := r.buf[r.r]
	if b == '\r' || b == '\n' {
		r.consumeNewline()
		return true
	}
	if r.cfg.Comment != 0 && rune(b) == r.cfg.Comment {
		r.consumeUntilNewline()
		return true
	}
	return false
}

// consumeUntilNewline skips entire comment lines until the line break.
func (r *Reader) consumeUntilNewline() {
	r.line++
	for r.r < r.w {
		b := r.buf[r.r]
		r.r++
		if b == '\n' {
			return
		}
	}
}
// trimLeadingWhitespace strips spaces if TrimLeadingSpace is enabled.
func (r *Reader) trimLeadingWhitespace() {
	if !r.cfg.TrimLeadingSpace {
		return
	}
	for r.r < r.w && (r.buf[r.r] == ' ' || r.buf[r.r] == '\t') {
		r.r++
	}
}
// scanUnquotedField parses bytes up to the delimiter or newline.
func (r *Reader) scanUnquotedField(delim byte) (bool, bool, error) {
	start := r.r
	for {
		if r.r >= r.w {
			if r.eof {
				r.scratch = append(r.scratch, r.buf[start:r.r]...)
				return false, true, nil
			}
			r.scratch = append(r.scratch, r.buf[start:r.r]...)
			if err := r.fill(); err != nil {
				return false, false, err
			}
			start = r.r
		}
		b := r.buf[r.r]
		if b == delim {
			r.scratch = append(r.scratch, r.buf[start:r.r]...)
			r.r++
			return false, false, nil
		}
		if b == '\r' || b == '\n' {
			r.scratch = append(r.scratch, r.buf[start:r.r]...)
			r.consumeNewline()
			return false, true, nil
		}
		r.r++
	}
}
// scanQuotedField parses RFC 4180 quoted bytes with double-quote escaping.
func (r *Reader) scanQuotedField(delim, quote byte) (bool, bool, error) {
	r.r++ // Skip opening quote
	hasEsc := false
	for {
		if r.r >= r.w {
			if r.eof {
				return false, false, newParseError(r.line, len(r.colOffs)+1, "unterminated quoted field", ErrQuote)
			}
			if err := r.fill(); err != nil {
				return false, false, err
			}
		}
		b := r.buf[r.r]
		if b == quote {
			r.r++
			if r.r < r.w && r.buf[r.r] == quote {
				r.scratch = append(r.scratch, quote)
				r.r++
				hasEsc = true
				continue
			}
			return r.consumeAfterQuote(delim, hasEsc)
		}
		r.scratch = append(r.scratch, b)
		r.r++
	}
}
// consumeAfterQuote processes trailing characters following the closing quote.
func (r *Reader) consumeAfterQuote(delim byte, hasEsc bool) (bool, bool, error) {
	for r.r < r.w && r.buf[r.r] != delim && r.buf[r.r] != '\r' && r.buf[r.r] != '\n' {
		r.r++
	}
	if r.r >= r.w && r.eof {
		return hasEsc, true, nil
	}
	if r.r < r.w && r.buf[r.r] == delim {
		r.r++
		return hasEsc, false, nil
	}
	r.consumeNewline()
	return hasEsc, true, nil
}
// parseRecord parses fields from the buffer into the current Record.
func (r *Reader) parseRecord() (bool, error) {
	r.colOffs = r.colOffs[:0]
	r.colLens = r.colLens[:0]
	r.scratch = r.scratch[:0]
	r.line++
	delim := byte(r.cfg.Delimiter)
	quote := byte(r.cfg.Quote)
	for {
		fieldStart := len(r.scratch)
		hasQuote, lastField, err := r.scanField(delim, quote)
		if err != nil {
			return false, err
		}
		fieldLen := len(r.scratch) - fieldStart
		r.colOffs = append(r.colOffs, fieldStart)
		r.colLens = append(r.colLens, fieldLen)
		if lastField {
			r.buildRecord(hasQuote)
			return true, nil
		}
	}
}

// scanField extracts a single column handling optional quotation and escaping.
func (r *Reader) scanField(delim, quote byte) (bool, bool, error) {
	r.trimLeadingWhitespace()
	if r.r < r.w && r.buf[r.r] == quote {
		return r.scanQuotedField(delim, quote)
	}
	return r.scanUnquotedField(delim)
}

// buildRecord updates the current record structure with parsed offsets.
func (r *Reader) buildRecord(hasEsc bool) {
	r.currRecord.raw = r.scratch
	r.currRecord.colOffs = r.colOffs
	r.currRecord.colLens = r.colLens
	r.currRecord.lineNum = r.line
	r.currRecord.numCols = len(r.colOffs)
	r.currRecord.hasEsc = hasEsc
	r.currRecord.err = nil
}