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