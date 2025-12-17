# Implementation Plan

- [x] 1. Set up project structure and CLI framework
  - Initialize Go module with `go mod init gobpftool`
  - Add dependencies: `github.com/cilium/ebpf` and `github.com/spf13/cobra`
  - Create directory structure: `cmd/`, `pkg/prog/`, `pkg/maps/`, `pkg/output/`, `internal/utils/`
  - Create `main.go` entry point
  - _Requirements: 1.1, 1.2_

- [x] 2. Implement root command with global flags
  - Create `cmd/root.go` with root command using cobra
  - Add global flags: `-j`/`--json`, `-p`/`--pretty`
  - Add `--version` flag support
  - Implement help output for root command
  - Write unit tests for flag parsing
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6_

- [x] 3. Implement output formatter
  - [x] 3.1 Create formatter interface in `pkg/output/formatter.go`
    - Define `Format` type and constants (Plain, JSON, JSONPretty)
    - Define `Formatter` interface with methods for programs, maps, entries
    - _Requirements: 9.1, 9.2, 9.3_

  - [x] 3.2 Implement plain text formatter in `pkg/output/plain.go`
    - Format programs matching bpftool output format
    - Format maps matching bpftool output format
    - Format map entries (dump, lookup, getnext)
    - Write unit tests with expected output strings
    - _Requirements: 9.3, 2.9_

  - [x] 3.3 Implement JSON formatter in `pkg/output/json.go`
    - Format programs as JSON with bpftool-compatible field names
    - Format maps as JSON
    - Format map entries as JSON
    - Support both compact and pretty-printed JSON
    - Write unit tests
    - _Requirements: 9.1, 9.2, 9.4_

- [x] 4. Implement program service
  - [x] 4.1 Create program service interface in `pkg/prog/service.go`
    - Define `ProgramInfo` struct with all required fields
    - Define `Service` interface with List, GetByID, GetByTag, GetByName, GetByPinnedPath
    - _Requirements: 2.1, 2.8_

  - [x] 4.2 Implement program service using cilium/ebpf in `pkg/prog/service_impl.go`
    - Implement `List()` using `ebpf.ProgramIterator`
    - Implement `GetByID()` using `ebpf.NewProgramFromID`
    - Implement `GetByTag()` by filtering list results
    - Implement `GetByName()` by filtering list results
    - Implement `GetByPinnedPath()` using `ebpf.LoadPinnedProgram`
    - Extract program info including map_ids from program info
    - Write unit tests (mock-based for interface, integration for real eBPF)
    - _Requirements: 2.1, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8_

- [ ] 5. Implement prog commands
  - [ ] 5.1 Create prog command structure in `cmd/prog.go`
    - Create `prog` parent command
    - Wire up to root command
    - _Requirements: 1.1_

  - [ ] 5.2 Implement `prog show` command
    - Parse program identifier (id, tag, name, pinned)
    - Call program service to get program info
    - Format output using formatter
    - Handle errors with appropriate messages
    - Support `list` alias
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [ ] 5.3 Implement `prog help` command
    - Display help for prog subcommands
    - Show available commands and options
    - _Requirements: 3.1, 3.2, 3.3_

- [ ] 6. Implement hex utilities
  - Create `internal/utils/hex.go` with hex parsing functions
  - Parse space-separated hex bytes to byte slice
  - Format byte slice to hex string output
  - Write unit tests for parsing and formatting
  - _Requirements: 6.3, 6.5_

- [ ] 7. Implement map service
  - [ ] 7.1 Create map service interface in `pkg/maps/service.go`
    - Define `MapInfo` struct with all required fields
    - Define `MapEntry` struct for key-value pairs
    - Define `Service` interface with List, GetByID, GetByName, GetByPinnedPath, Dump, Lookup, GetNextKey
    - _Requirements: 4.1, 4.7_

  - [ ] 7.2 Implement map service using cilium/ebpf in `pkg/maps/service_impl.go`
    - Implement `List()` using `ebpf.MapIterator`
    - Implement `GetByID()` using `ebpf.NewMapFromID`
    - Implement `GetByName()` by filtering list results
    - Implement `GetByPinnedPath()` using `ebpf.LoadPinnedMap`
    - Implement `Dump()` by iterating all entries
    - Implement `Lookup()` using map's Lookup method
    - Implement `GetNextKey()` using map's NextKey method
    - Write unit tests
    - _Requirements: 4.1, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 8. Implement map commands
  - [ ] 8.1 Create map command structure in `cmd/map.go`
    - Create `map` parent command
    - Wire up to root command
    - _Requirements: 1.1_

  - [ ] 8.2 Implement `map show` command
    - Parse map identifier (id, name, pinned)
    - Call map service to get map info
    - Format output using formatter
    - Handle errors with appropriate messages
    - Support `list` alias
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_

  - [ ] 8.3 Implement `map dump` command
    - Parse map identifier
    - Call map service Dump method
    - Format all entries using formatter
    - Display element count
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_

  - [ ] 8.4 Implement `map lookup` command
    - Parse map identifier and key data
    - Parse hex key bytes using utils
    - Call map service Lookup method
    - Format result using formatter
    - Handle key not found error
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [ ] 8.5 Implement `map getnext` command
    - Parse map identifier and optional key data
    - Call map service GetNextKey method
    - Format result showing current and next key
    - Handle no more keys case
    - Handle empty map case
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [ ] 8.6 Implement `map help` command
    - Display help for map subcommands
    - Show available commands and options
    - _Requirements: 8.1, 8.2, 8.3_

- [ ] 9. Implement error handling
  - Create error types in `pkg/errors/errors.go`
  - Implement permission error detection and messaging
  - Implement BPF filesystem mount detection
  - Ensure proper exit codes (0 for success, 1 for failure)
  - Write error messages to stderr
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ] 10. Add version command
  - Create `cmd/version.go`
  - Display version information
  - Support `--version` flag on root command
  - _Requirements: 1.3_
