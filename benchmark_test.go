package rumil_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"testing"

	"github.com/wesleyskap/rumil-csv"
)

func generateBenchmarkCSV(rows int) []byte {
	var buf bytes.Buffer
	for i := 0; i < rows; i++ {
		buf.WriteString("1001,John Doe,Engineering,San Francisco,95000.50,true,\"Senior Staff Software Engineer, Lead\"\n")
	}
	return buf.Bytes()
}

func BenchmarkRumilReader(b *testing.B) {
	data := generateBenchmarkCSV(1000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := rumil.NewReader(bytes.NewReader(data))
		count := 0
		for r.Scan() {
			rec := r.Record()
			_ = rec.At(0)
			count++
		}
		if count != 1000 {
			b.Fatalf("expected 1000 rows, got %d", count)
		}
	}
}

func BenchmarkStdCSVReader(b *testing.B) {
	data := generateBenchmarkCSV(1000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := csv.NewReader(bytes.NewReader(data))
		count := 0
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("std csv error: %v", err)
			}
			_ = record[0]
			count++
		}
		if count != 1000 {
			b.Fatalf("expected 1000 rows, got %d", count)
		}
	}
}

func BenchmarkRumilWriter(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := rumil.NewWriter(io.Discard)
		for j := 0; j < 1000; j++ {
			_ = w.WriteStringRow("1001", "John Doe", "Engineering", "San Francisco", "95000.50", "true")
		}
		_ = w.Flush()
	}
}