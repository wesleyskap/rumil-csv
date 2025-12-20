package rumil

import (
	"bytes"
	"io"
)

// Writer writes delimiter-separated data to an output stream using buffered writes.
// Struct fields are arranged from largest byte size to smallest for optimal memory layout.
type Writer struct {
	buf []byte       // 24 bytes
	out io.Writer    // 16 bytes
	err error        // 16 bytes
	cfg WriterConfig // 16 bytes (ordered struct)
}

// NewWriter constructs a high-throughput buffered CSV writer.
func NewWriter(w io.Writer, opts ...WriterOption) *Writer {
	cfg := defaultWriterConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Writer{
		out: w,
		buf: make([]byte, 0, cfg.BufferSize),
		cfg: cfg,
	}
}

// Reset resets the internal buffer and re-targets the writer to a new destination.
func (w *Writer) Reset(out io.Writer) {
	w.out = out
	w.buf = w.buf[:0]
	w.err = nil
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	if w.err != nil {
		return w.err
	}
	if len(w.buf) == 0 {
		return nil
	}
	_, err := w.out.Write(w.buf)
	w.buf = w.buf[:0]
	w.err = err
	return err
}
// appendField formats and appends a single field, applying RFC 4180 quotes when necessary.
func (w *Writer) appendField(field []byte, delim, quote byte) {
	if !w.cfg.AlwaysQuote && !w.needsQuotes(field, delim, quote) {
		w.buf = append(w.buf, field...)
		return
	}
	w.buf = append(w.buf, quote)
	for _, b := range field {
		if b == quote {
			w.buf = append(w.buf, quote, quote)
		} else {
			w.buf = append(w.buf, b)
		}
	}
	w.buf = append(w.buf, quote)
}

// needsQuotes checks whether a field contains delimiter, quote, or newline characters.
func (w *Writer) needsQuotes(field []byte, delim, quote byte) bool {
	if len(field) == 0 {
		return false
	}
	return bytes.IndexByte(field, delim) >= 0 ||
		bytes.IndexByte(field, quote) >= 0 ||
		bytes.IndexByte(field, '\r') >= 0 ||
		bytes.IndexByte(field, '\n') >= 0
}