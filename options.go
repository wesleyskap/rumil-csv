package rumil

// ReaderConfig holds configuration parameters for CSV stream parsing.
// Fields are ordered from largest byte size to smallest to optimize memory layout.
type ReaderConfig struct {
	BufferSize       int  // 8 bytes
	MaxRecordSize    int  // 8 bytes
	Delimiter        rune // 4 bytes
	Quote            rune // 4 bytes
	Comment          rune // 4 bytes
	TrimLeadingSpace bool // 1 byte
	LazyQuotes       bool // 1 byte
	ReuseRecord      bool // 1 byte
}

// Option represents a functional option for configuring a Reader.
type Option func(*ReaderConfig)

// defaultReaderConfig creates a default configuration matching RFC 4180.
func defaultReaderConfig() ReaderConfig {
	return ReaderConfig{
		BufferSize:       64 * 1024,
		MaxRecordSize:    10 * 1024 * 1024,
		Delimiter:        ',',
		Quote:            '"',
		Comment:          0,
		TrimLeadingSpace: false,
		LazyQuotes:       false,
		ReuseRecord:      true,
	}
}