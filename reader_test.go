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
func TestReaderQuotedFieldsAndEscaping(t *testing.T) {
	input := "header1,header2\n\"val,with,comma\",\"val \"\"with\"\" quotes\"\n"
	r := rumil.NewReader(strings.NewReader(input))

	if !r.Scan() {
		t.Fatalf("failed scanning header: %v", r.Err())
	}
	if !r.Scan() {
		t.Fatalf("failed scanning data: %v", r.Err())
	}

	rec := r.Record()
	if rec.StringAt(0) != "val,with,comma" {
		t.Fatalf("expected comma content, got %q", rec.StringAt(0))
	}
	if rec.StringAt(1) != "val \"with\" quotes" {
		t.Fatalf("expected escaped quote content, got %q", rec.StringAt(1))
	}
}

func TestReaderMultilineQuotes(t *testing.T) {
	input := "a,b\n\"line 1\nline 2\",c\n"
	r := rumil.NewReader(strings.NewReader(input))

	if !r.Scan() {
		t.Fatalf("failed scanning header: %v", r.Err())
	}
	if !r.Scan() {
		t.Fatalf("failed scanning multiline data: %v", r.Err())
	}

	rec := r.Record()
	expected := "line 1\nline 2"
	if rec.StringAt(0) != expected {
		t.Fatalf("expected multiline %q, got %q", expected, rec.StringAt(0))
	}
	if rec.StringAt(1) != "c" {
		t.Fatalf("expected c, got %q", rec.StringAt(1))
	}
}
func TestReaderCustomDelimiterAndComments(t *testing.T) {
	input := "# Leading comment\nitem;qty;price;active\nSword;10;120.50;true\n# Trailing comment\nShield;5;85.00;false\n"
	r := rumil.NewReader(
		strings.NewReader(input),
		rumil.WithDelimiter(';'),
		rumil.WithComment('#'),
	)

	if !r.Scan() {
		t.Fatalf("failed scanning header: %v", r.Err())
	}
	if !r.Scan() {
		t.Fatalf("failed scanning row 1: %v", r.Err())
	}
	rec1 := r.Record()
	price, err := rec1.FloatAt(2)
	if err != nil || price != 120.50 {
		t.Fatalf("expected price 120.50, got %f, err: %v", price, err)
	}
	active, err := rec1.BoolAt(3)
	if err != nil || !active {
		t.Fatalf("expected active true, got %v, err: %v", active, err)
	}

	if !r.Scan() {
		t.Fatalf("failed scanning row 2: %v", r.Err())
	}
	rec2 := r.Record()
	if rec2.StringAt(0) != "Shield" {
		t.Fatalf("expected Shield, got %s", rec2.StringAt(0))
	}

	if r.Scan() {
		t.Fatalf("unexpected extra row")
	}
}