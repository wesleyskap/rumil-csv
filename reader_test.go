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
func TestReaderSimpleUnquoted(t *testing.T) {
	input := "id,name,age\n1,Alice,30\n2,Bob,25\n"
	r := rumil.NewReader(strings.NewReader(input))

	if !r.Scan() {
		t.Fatalf("expected first row, got none: %v", r.Err())
	}
	rec := r.Record()
	if rec.Len() != 3 || rec.StringAt(0) != "id" || rec.StringAt(1) != "name" || rec.StringAt(2) != "age" {
		t.Fatalf("unexpected header record: %+v", rec.Strings())
	}

	if !r.Scan() {
		t.Fatalf("expected second row, got none: %v", r.Err())
	}
	id, err := r.Record().IntAt(0)
	if err != nil || id != 1 {
		t.Fatalf("expected id 1, got %d, err %v", id, err)
	}
	age, err := r.Record().IntAt(2)
	if err != nil || age != 30 {
		t.Fatalf("expected age 30, got %d, err %v", age, err)
	}

	if !r.Scan() {
		t.Fatalf("expected third row, got none: %v", r.Err())
	}
	if r.Record().StringAt(1) != "Bob" {
		t.Fatalf("expected Bob, got %s", r.Record().StringAt(1))
	}

	if r.Scan() {
		t.Fatalf("unexpected extra record after EOF")
	}
	if r.Err() != nil {
		t.Fatalf("unexpected error: %v", r.Err())
	}
}