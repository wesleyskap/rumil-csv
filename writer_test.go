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