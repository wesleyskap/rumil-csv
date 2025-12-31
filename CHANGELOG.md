# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-31

### Added
- Zero-allocation streaming CSV reader core engine.
- Memory-aligned Record structure ordered by field size.
- Range-over-func iterator support (iter.Seq2) for Go 1.23+.
- Classic Scan and Record methods compatible with standard scanning workflows.
- RFC 4180 strict parsing with multi-line fields and escaped quotes.
- Typed zero-copy scalar converters (IntAt, FloatAt, BoolAt).
- Functional configuration options for delimiter, quote, buffer sizing, and comments.
- High-throughput buffered Writer with automatic RFC 4180 escaping.
- Benchmark suite demonstrating 0 B/op and 0 allocs/op during record streaming.
- Comprehensive unit test coverage for reader, writer, and error states.