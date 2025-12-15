package rumil

import (
	"iter"
)

// All yields an iter.Seq2 sequence for range-over-func iteration introduced in Go 1.23.
func (r *Reader) All() iter.Seq2[*Record, error] {
	return func(yield func(*Record, error) bool) {
		for r.Scan() {
			if !yield(r.Record(), nil) {
				return
			}
		}
		if err := r.Err(); err != nil {
			yield(nil, err)
		}
	}
}