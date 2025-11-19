# gocachectl

[![CI](https://github.com/muhammadali7768/gocachectl/workflows/CI/badge.svg)](https://github.com/muhammadali7768/gocachectl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhammadali7768/gocachectl)](https://goreportcard.com/report/github.com/muhammadali7768/gocachectl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Universal Go Cache Management Tool** - Manage all your Go caches from a single CLI.

## Why gocachectl?

Go maintains multiple caches across your system:
- **Build Cache** (`GOCACHE`) - Compiled packages and build artifacts
- **Module Cache** (`GOMODCACHE`) - Downloaded dependencies
- **Test Cache** - Cached test results

There's no unified way to view, analyze, or manage these caches. `gocachectl` solves this problem.

## Features

-  **Unified Statistics** - View all cache sizes and counts at once
-  **Cache Information** - See cache locations and Go environment
-  **Selective Clearing** - Clear all caches or specific ones
-  **Space Management** - Know exactly how much disk space your caches use
-  **Multiple Output Formats** - Human-readable tables or JSON



## Usage

### View Cache Statistics

```bash
# Show all cache statistics
gocachectl stats

# Show only build cache
gocachectl stats --build

# Show only module cache
gocachectl stats --modules

# Show only test cache
gocachectl stats --test

# JSON output
gocachectl stats --json

# Verbose output with more details
gocachectl stats --verbose
```



### Show Cache Information

```bash
# Show cache locations and Go version
gocachectl info

# JSON output
gocachectl info --json
```

### Clear Caches

```bash
# Clear all caches (with confirmation)
gocachectl clear --all

# Clear all without confirmation
gocachectl clear --all --force

# Clear only build cache
gocachectl clear --build

# Clear only module cache
gocachectl clear --modules

# Clear only test cache
gocachectl clear --test

# Dry run - see what would be deleted
gocachectl clear --all --dry-run

# Quiet mode (minimal output)
gocachectl clear --all --force --quiet
```

### Show Version

```bash
gocachectl version

# Verbose version info
gocachectl version --verbose
```

## Global Flags

Available for all commands:

- `--verbose`, `-v` - Enable verbose output
- `--json` - Output in JSON format
- `--quiet`, `-q` - Minimal output (errors only)
- `--help`, `-h` - Show help message
