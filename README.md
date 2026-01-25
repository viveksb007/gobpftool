# gobpftool

[![Build](https://github.com/viveksb007/gobpftool/actions/workflows/build.yml/badge.svg)](https://github.com/viveksb007/gobpftool/actions)

A Go-based CLI tool for inspecting eBPF programs and maps, similar to Linux `bpftool`.

## Installation

### Homebrew (Linux)

```bash
brew tap viveksb007/tap
brew install gobpftool
```

### Download Binary

Download the latest release from the [releases page](https://github.com/viveksb007/gobpftool/releases).

### Build from source

```bash
make build
```

## Usage

Most commands require root privileges to access eBPF subsystem.

### Program Commands

![Program Show](docs/prog_show.png)



```bash
# List all loaded programs
sudo ./gobpftool prog show

# Show program by ID
sudo ./gobpftool prog show id 123

# Show programs by name
sudo ./gobpftool prog show name my_prog

# Show programs by tag
sudo ./gobpftool prog show tag f0055c08993fea1e

# Show pinned program
sudo ./gobpftool prog show pinned /sys/fs/bpf/my_prog
```

### Map Commands

![Map Show](docs/map_show.png)

```bash
# List all loaded maps
sudo ./gobpftool map show

# Show map by ID
sudo ./gobpftool map show id 123

# Show maps by name
sudo ./gobpftool map show name my_map

# Show pinned map
sudo ./gobpftool map show pinned /sys/fs/bpf/my_map

# Dump all entries in a map
sudo ./gobpftool map dump id 123

# Lookup a key (hex bytes)
sudo ./gobpftool map lookup id 123 key 00 00 00 00

# Get first key
sudo ./gobpftool map getnext id 123

# Get next key after specified key
sudo ./gobpftool map getnext id 123 key 00 00 00 00
```

### Output Formats

```bash
# JSON output
sudo ./gobpftool -j map show

# Pretty-printed JSON
sudo ./gobpftool -p prog show
```

## License

MIT
