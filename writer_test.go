package rumil_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wesleyskap/rumil-csv"
)

func TestWriterBasicRows(t *testing.T) {
	var buf bytes.Buffer
	w := rumil.NewWriter(&buf)

	err := w.WriteStringRow("id", "name", "role")
	if err != nil {
		t.Fatalf("failed writing header: %v", err)
	}
	err = w.WriteStringRow("1", "Gandalf", "Wizard")
	if err != nil {
		t.Fatalf("failed writing row: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("failed to flush: %v", err)
	}

	expected := "id,name,role\r\n1,Gandalf,Wizard\r\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output: %q, expected: %q", buf.String(), expected)
	}
}

func TestWriterQuotingTriggers(t *testing.T) {
	var buf bytes.Buffer
	w := rumil.NewWriter(&buf)

	err := w.WriteStringRow("normal", "field,with,comma", "field \"with\" quotes", "field\nwith\nnewlines")
	if err != nil {
		t.Fatalf("failed writing row: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("failed flushing: %v", err)
	}

	expected := "normal,\"field,with,comma\",\"field \"\"with\"\" quotes\",\"field\nwith\nnewlines\"\r\n"
	if buf.String() != expected {
		t.Fatalf("unexpected escaped output: %q, expected %q", buf.String(), expected)
	}
}
func TestWriterAlwaysQuote(t *testing.T) {
	var buf bytes.Buffer
	w := rumil.NewWriter(&buf, rumil.WithAlwaysQuote(true))

	err := w.WriteStringRow("a", "b", "c")
	if err != nil {
		t.Fatalf("failed writing: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("failed flush: %v", err)
	}

	expected := "\"a\",\"b\",\"c\"\r\n"
	if buf.String() != expected {
		t.Fatalf("expected always quoted: %q, got: %q", expected, buf.String())
	}
}

func TestWriterWriteAll(t *testing.T) {
	var buf bytes.Buffer
	w := rumil.NewWriter(&buf)

	rows := [][]string{
		{"num", "label"},
		{"100", "first"},
		{"200", "second"},
	}

	if err := w.WriteStringAll(rows); err != nil {
		t.Fatalf("failed WriteStringAll: %v", err)
	}

	expected := "num,label\r\n100,first\r\n200,second\r\n"
	if buf.String() != expected {
		t.Fatalf("unexpected output: %q", buf.String())
	}

	r := rumil.NewReader(strings.NewReader(buf.String()))
	rowCount := 0
	for r.Scan() {
		rowCount++
	}
	if rowCount != 3 {
		t.Fatalf("expected 3 roundtrip rows, read %d", rowCount)
	}
}