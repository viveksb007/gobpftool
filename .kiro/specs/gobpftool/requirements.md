# Requirements Document

## Introduction

gobpftool is a Go-based CLI tool that provides functionality similar to the Linux `bpftool` utility for inspecting eBPF programs and maps. It leverages the [cilium/ebpf](https://github.com/cilium/ebpf) library to interact with the kernel's eBPF subsystem.

**Scope for Initial Implementation:**
- Program inspection: `prog show`, `prog help`
- Map read operations: `map show`, `map dump`, `map lookup`, `map getnext`, `map help`

Future iterations may add program loading, map creation, and mutating operations.

## Requirements

### Requirement 1: CLI Framework and Basic Structure

**User Story:** As a developer, I want a well-structured CLI with subcommands for `prog` and `map`, so that I can interact with eBPF resources using familiar bpftool-like syntax.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool` without arguments THEN the system SHALL display help information showing available subcommands (`prog`, `map`).
2. WHEN the user runs `gobpftool --help` THEN the system SHALL display detailed usage information including all available subcommands and global options.
3. WHEN the user runs `gobpftool --version` THEN the system SHALL display the version information of the tool.
4. WHEN the user runs `gobpftool` with an invalid subcommand THEN the system SHALL display an error message and suggest valid subcommands.
5. WHEN the user runs any command THEN the system SHALL support `-j` or `--json` flag to output results in JSON format.
6. WHEN the user runs any command THEN the system SHALL support `-p` or `--pretty` flag to output results in pretty-printed JSON format.

### Requirement 2: Program Listing (prog show/list)

**User Story:** As a system administrator, I want to list all loaded eBPF programs, so that I can see what programs are currently active in the kernel.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool prog show` THEN the system SHALL display a list of all loaded eBPF programs with their ID, type, name, and tag.
2. WHEN the user runs `gobpftool prog list` THEN the system SHALL behave identically to `prog show`.
3. WHEN the user runs `gobpftool prog show id <ID>` THEN the system SHALL display detailed information about the specific program with that ID.
4. WHEN the user runs `gobpftool prog show tag <TAG>` THEN the system SHALL display information about programs matching that tag.
5. WHEN the user runs `gobpftool prog show pinned <PATH>` THEN the system SHALL display information about the program pinned at that BPF filesystem path.
6. WHEN the user runs `gobpftool prog show name <NAME>` THEN the system SHALL display information about programs matching that name.
7. IF no eBPF programs are loaded THEN the system SHALL display an empty list or appropriate message.
8. WHEN displaying program information THEN the system SHALL include: program ID, type, name, tag, GPL compatible flag, loaded time (with timezone), uid, bytes translated (xlated), bytes JITed (jited), memory lock size (memlock), and associated map IDs (map_ids).
9. WHEN displaying program information THEN the output format SHALL match bpftool's format:
   ```
   <ID>: <type>  name <name>  tag <tag>  gpl
           loaded_at <timestamp>  uid <uid>
           xlated <bytes>B  jited <bytes>B  memlock <bytes>B  map_ids <id1>,<id2>,...
   ```

### Requirement 3: Program Help (prog help)

**User Story:** As a user, I want built-in help for program commands, so that I can quickly reference available options and usage.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool prog help` THEN the system SHALL display help information specific to program commands.
2. WHEN the user runs `gobpftool prog --help` THEN the system SHALL display help for program subcommands.
3. WHEN displaying help THEN the system SHALL show command syntax, available options, and brief descriptions.

### Requirement 4: Map Listing (map show/list)

**User Story:** As a system administrator, I want to list all eBPF maps, so that I can see what maps are currently loaded in the kernel.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool map show` THEN the system SHALL display a list of all loaded eBPF maps with their ID, type, name, key size, value size, and max entries.
2. WHEN the user runs `gobpftool map list` THEN the system SHALL behave identically to `map show`.
3. WHEN the user runs `gobpftool map show id <ID>` THEN the system SHALL display detailed information about the specific map with that ID.
4. WHEN the user runs `gobpftool map show pinned <PATH>` THEN the system SHALL display information about the map pinned at that BPF filesystem path.
5. WHEN the user runs `gobpftool map show name <NAME>` THEN the system SHALL display information about maps matching that name.
6. IF no eBPF maps are loaded THEN the system SHALL display an empty list or appropriate message.
7. WHEN displaying map information THEN the system SHALL include: map ID, type, name, key size, value size, max entries, flags, and memory usage.

### Requirement 5: Map Dump (map dump)

**User Story:** As a developer, I want to dump all entries in an eBPF map, so that I can inspect the complete contents of a map.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool map dump id <ID>` THEN the system SHALL display all key-value pairs in the map.
2. WHEN the user runs `gobpftool map dump pinned <PATH>` THEN the system SHALL display all key-value pairs in the pinned map.
3. WHEN the user runs `gobpftool map dump name <NAME>` THEN the system SHALL display all key-value pairs in maps matching that name.
4. WHEN dumping map contents THEN the system SHALL display keys and values in hex format.
5. IF the map is empty THEN the system SHALL display "Found 0 elements" or similar message.
6. WHEN dump completes THEN the system SHALL display the total number of elements found.

### Requirement 6: Map Lookup (map lookup)

**User Story:** As a developer, I want to look up specific entries in an eBPF map, so that I can inspect individual values.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool map lookup id <ID> key <KEY_DATA>` THEN the system SHALL display the value associated with the specified key.
2. WHEN the user runs `gobpftool map lookup pinned <PATH> key <KEY_DATA>` THEN the system SHALL display the value for the key in the pinned map.
3. WHEN specifying key data THEN the system SHALL support space-separated hex bytes (e.g., `key 0a 0b 0c 0d`).
4. IF the key does not exist THEN the system SHALL display an appropriate error message.
5. WHEN displaying lookup results THEN the system SHALL show both the key and value in hex format.

### Requirement 7: Map Get Next Key (map getnext)

**User Story:** As a developer, I want to iterate through map keys, so that I can traverse all entries in a map.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool map getnext id <ID>` THEN the system SHALL display the first key in the map.
2. WHEN the user runs `gobpftool map getnext id <ID> key <KEY_DATA>` THEN the system SHALL display the next key after the specified key.
3. WHEN the user runs `gobpftool map getnext pinned <PATH>` THEN the system SHALL display the first key in the pinned map.
4. IF there are no more keys THEN the system SHALL indicate that iteration is complete.
5. IF the map is empty THEN the system SHALL display an appropriate message.

### Requirement 8: Map Help (map help)

**User Story:** As a user, I want built-in help for map commands, so that I can quickly reference available options and usage.

#### Acceptance Criteria

1. WHEN the user runs `gobpftool map help` THEN the system SHALL display help information specific to map commands.
2. WHEN the user runs `gobpftool map --help` THEN the system SHALL display help for map subcommands.
3. WHEN displaying help THEN the system SHALL show command syntax, available options, and brief descriptions.

### Requirement 9: Output Formatting

**User Story:** As a user, I want flexible output formats, so that I can integrate gobpftool with scripts and other tools.

#### Acceptance Criteria

1. WHEN the user specifies `-j` or `--json` flag THEN the system SHALL output results in compact JSON format.
2. WHEN the user specifies `-p` or `--pretty` flag THEN the system SHALL output results in pretty-printed JSON format with indentation.
3. WHEN no format flag is specified THEN the system SHALL output results in human-readable plain text format.
4. WHEN outputting in JSON format THEN the system SHALL use consistent field names matching bpftool's JSON output where applicable.

### Requirement 10: Error Handling and Permissions

**User Story:** As a user, I want clear error messages when operations fail, so that I can understand and resolve issues.

#### Acceptance Criteria

1. IF the user lacks CAP_SYS_ADMIN or CAP_BPF capabilities THEN the system SHALL display a clear error message about insufficient permissions.
2. IF the BPF filesystem is not mounted THEN the system SHALL display instructions on how to mount it.
3. WHEN any operation fails THEN the system SHALL return a non-zero exit code.
4. WHEN any operation succeeds THEN the system SHALL return exit code 0.
5. WHEN an error occurs THEN the system SHALL display a descriptive error message to stderr.
