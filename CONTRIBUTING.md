# Contributing to rumil-csv

Thank you for your interest in contributing to rumil-csv.

## Development Standards

1. Code Quality
   - All code must follow standard Go conventions formatted with gofmt.
   - Run go vet ./... and ensure zero warnings or lint errors.
   - All struct definitions must order fields from largest byte size to smallest to prevent unnecessary padding.
   - Maintain function lengths between 4 and 20 lines.

2. Testing and Benchmarks
   - Every bug fix or new feature must include unit tests.
   - Any modifications to reader or writer internals must verify zero heap allocations using go test -bench=. -benchmem.

3. Git Commits and Messages
   - Use Conventional Commits formatting (feat:, fix:, test:, refactor:, docs:, chore:).
   - Write all commit messages and comments in clear English.
   - Do not use emojis in commit messages or code files.