package rumil_test

import (
	"strings"
	"testing"

	"github.com/wesleyskap/rumil-csv"
)

func TestReaderRecordClone(t *testing.T) {
	input := "alpha,beta\ngamma,delta\n"
	r := rumil.NewReader(strings.NewReader(input))

	if !r.Scan() {
		t.Fatalf("failed to scan row 1")
	}
	cloned := r.Record().Clone()

	if !r.Scan() {
		t.Fatalf("failed to scan row 2")
	}

	if cloned.StringAt(0) != "alpha" || cloned.StringAt(1) != "beta" {
		t.Fatalf("cloned record mutated: %v", cloned.Strings())
	}
	if r.Record().StringAt(0) != "gamma" || r.Record().StringAt(1) != "delta" {
		t.Fatalf("current record unexpected: %v", r.Record().Strings())
	}
}