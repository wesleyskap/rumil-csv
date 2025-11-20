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
// IntAt parses column index i as a 64-bit signed integer.
func (r *Record) IntAt(i int) (int64, error) {
	b := r.At(i)
	if b == nil {
		return 0, newParseError(r.lineNum, i+1, "column index out of bounds", nil)
	}
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, newParseError(r.lineNum, i+1, "failed to parse integer", err)
	}
	return v, nil
}
// FloatAt parses column index i as a 64-bit floating point number.
func (r *Record) FloatAt(i int) (float64, error) {
	b := r.At(i)
	if b == nil {
		return 0, newParseError(r.lineNum, i+1, "column index out of bounds", nil)
	}
	v, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, newParseError(r.lineNum, i+1, "failed to parse float", err)
	}
	return v, nil
}
// BoolAt parses column index i as a boolean value.
func (r *Record) BoolAt(i int) (bool, error) {
	b := r.At(i)
	if b == nil {
		return false, newParseError(r.lineNum, i+1, "column index out of bounds", nil)
	}
	v, err := strconv.ParseBool(string(b))
	if err != nil {
		return false, newParseError(r.lineNum, i+1, "failed to parse boolean", err)
	}
	return v, nil
}
// Fields returns all columns as a slice of byte slices referencing the internal buffer.
func (r *Record) Fields() [][]byte {
	if r == nil || r.numCols == 0 {
		return nil
	}
	out := make([][]byte, r.numCols)
	for i := 0; i < r.numCols; i++ {
		out[i] = r.At(i)
	}
	return out
}

// Strings allocates and returns all columns as Go strings.
func (r *Record) Strings() []string {
	if r == nil || r.numCols == 0 {
		return nil
	}
	out := make([]string, r.numCols)
	for i := 0; i < r.numCols; i++ {
		out[i] = r.StringAt(i)
	}
	return out
}
// Clone creates an independent heap copy of the record with isolated byte slices.
func (r *Record) Clone() *Record {
	if r == nil {
		return nil
	}
	rawCopy := make([]byte, len(r.raw))
	copy(rawCopy, r.raw)
	offsCopy := make([]int, r.numCols)
	lensCopy := make([]int, r.numCols)
	copy(offsCopy, r.colOffs[:r.numCols])
	copy(lensCopy, r.colLens[:r.numCols])
	return &Record{
		raw:     rawCopy,
		colOffs: offsCopy,
		colLens: lensCopy,
		err:     r.err,
		lineNum: r.lineNum,
		numCols: r.numCols,
		hasEsc:  r.hasEsc,
	}
}