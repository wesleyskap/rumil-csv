# rumil-csv

High-performance, zero-allocation streaming CSV and DSV toolkit for modern Go (1.23+).

## Overview

`rumil-csv` is a zero-dependency Go library engineered to dissect and write delimiter-separated values with zero allocations per record. By leveraging memory-aligned data structures, reusable sliding buffers, and Go 1.23 native range-over-func iterators, `rumil-csv` delivers maximum throughput while eliminating garbage collection pressure.

## Features

- Zero allocations per scanned record (0 B/op, 0 allocs/op in streaming mode).
- Go 1.23+ native iterator support via iter.Seq2 for idiomatic range-over-func loops.
- Memory-aligned structures optimized from largest to smallest byte size to avoid CPU cache padding.
- Strict RFC 4180 compliance including multi-line fields and escaped quotation marks.
- Custom delimiter and delimiter-separated format support (TSV, semicolon, pipe).
- Typed zero-copy scalar converters (IntAt, FloatAt, BoolAt).
- High-throughput buffered writer with automated RFC 4180 quoting and custom buffering.
- Zero external dependencies. Standard library only.

## Performance Benchmarks

Benchmarks executed on AMD Ryzen 7 5700X (Go 1.26 windows/amd64):

| Operation | Benchmark | Speed (ns/op) | Memory (B/op) | Allocations (allocs/op) |
| :--- | :--- | :--- | :--- | :--- |
| **Record Streaming** | `BenchmarkRumilScanRecord` | **232.9 ns/op** | **0 B/op** | **0 allocs/op** |
| 1,000 Rows Scan | `BenchmarkRumilReader` | 199.5 us/op | 132 KB/op | 7 allocs/op |
| 1,000 Rows Scan | `BenchmarkStdCSVReader` | 181.7 us/op | 212 KB/op | 2,016 allocs/op |
| 1,000 Rows Write | `BenchmarkRumilWriter` | 130.6 us/op | 147 KB/op | 5 allocs/op |

When scanning records continuously, `rumil-csv` achieves true 0 B/op and 0 allocs/op, compared to 2 allocations per row in Go standard `encoding/csv`.

## Installation

```bash
go get github.com/wesleyskap/rumil-csv
```

Requires Go 1.23 or newer.

## Quick Start

### Modern Go 1.23 Range-Over-Func Iteration

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wesleyskap/rumil-csv"
)

func main() {
	input := "id,name,role\n1,Feanor,Craftsman\n2,Rumil,Loremaster\n"
	reader := rumil.NewReader(strings.NewReader(input))

	for record, err := range reader.All() {
		if err != nil {
			panic(err)
		}
		id, _ := record.IntAt(0)
		name := record.StringAt(1)
		role := record.StringAt(2)
		fmt.Printf("Record %d: %s (%s)\n", id, name, role)
	}
}
```

### Classic Scanner Loop

```go
reader := rumil.NewReader(file, rumil.WithDelimiter(';'))

for reader.Scan() {
	rec := reader.Record()
	price, err := rec.FloatAt(2)
	if err != nil {
		log.Printf("Invalid price on line %d: %v", rec.Line(), err)
		continue
	}
	processItem(rec.At(0), price)
}

if err := reader.Err(); err != nil {
	log.Fatalf("Parse error: %v", err)
}
```

### High-Throughput Buffered Writer

```go
package main

import (
	"os"

	"github.com/wesleyskap/rumil-csv"
)

func main() {
	file, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := rumil.NewWriter(file)
	defer writer.Flush()

	writer.WriteStringRow("id", "title", "description")
	writer.WriteStringRow("1", "Sarati Script", "Created by Rumil, precursor to Tengwar")
	writer.WriteStringRow("2", "Silmarils", "Crafted by Feanor, containing Valinor light")
}
```

## Configuration Options

### Reader Options

- `WithDelimiter(r rune)`: Sets custom field delimiter (default `,`).
- `WithQuote(r rune)`: Sets custom quote rune (default `"`).
- `WithComment(r rune)`: Skips lines starting with specified rune.
- `WithTrimLeadingSpace(trim bool)`: Strips leading spaces on unquoted fields.
- `WithBufferSize(size int)`: Configures internal read buffer (default 64KB).
- `WithMaxRecordSize(size int)`: Sets upper threshold for growing record buffers.
- `WithReuseRecord(reuse bool)`: Reuses underlying slice memory across scan invocations.

### Writer Options

- `WithWriterDelimiter(r rune)`: Sets delimiter rune for output formatting.
- `WithWriterQuote(r rune)`: Sets quote character for escaping.
- `WithAlwaysQuote(always bool)`: Quotes all output fields unconditionally.

## License

MIT License. Copyright (c) 2025-2026 Wesley Skap.