package rumil

import (
	"errors"
	"fmt"
)

// Sentinel errors for parser state and operations.
var (
	ErrQuote          = errors.New("rumil: excess or unclosed quote in field")
	ErrFieldCount      = errors.New("rumil: unexpected number of fields in record")
	ErrBufferExceeded = errors.New("rumil: record exceeds internal buffer capacity")
	ErrClosed         = errors.New("rumil: reader or writer is closed")
)

// ParseError records a syntax or structure failure during CSV scanning.
// Struct fields are ordered largest to smallest byte size for optimal memory alignment.
type ParseError struct {
	Message string // 16 bytes (pointer + int)
	Err     error  // 16 bytes (interface: type + data ptr)
	Line    int64  // 8 bytes
	Column  int    // 8 bytes (on 64-bit platforms)
}