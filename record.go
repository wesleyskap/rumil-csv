package rumil

import (
	"strconv"
)

// Record represents a single parsed CSV row fanning out over zero-copy byte slices.
// Fields are ordered from largest byte size to smallest to prevent struct memory padding.
type Record struct {
	raw     []byte // 24 bytes
	colOffs []int  // 24 bytes
	colLens []int  // 24 bytes
	err     error  // 16 bytes
	lineNum int64  // 8 bytes
	numCols int    // 8 bytes
	hasEsc  bool   // 1 byte
}
// Len returns the count of columns detected in the current record.
func (r *Record) Len() int {
	if r == nil {
		return 0
	}
	return r.numCols
}

// Line returns the source input line number for this record.
func (r *Record) Line() int64 {
	if r == nil {
		return 0
	}
	return r.lineNum
}
// At returns the raw byte slice for column index i without allocating heap memory.
func (r *Record) At(i int) []byte {
	if r == nil || i < 0 || i >= r.numCols {
		return nil
	}
	start := r.colOffs[i]
	end := start + r.colLens[i]
	return r.raw[start:end]
}
// StringAt returns the column content as a standard Go string.
func (r *Record) StringAt(i int) string {
	b := r.At(i)
	if b == nil {
		return ""
	}
	return string(b)
}