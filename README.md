# CDN IP Scanner - Clean IP Discovery Tool

[![Xray-Core](https://img.shields.io/badge/Xray--Core-v26.1.23-blue)](https://github.com/XTLS/Xray-core)
[![Go](https://img.shields.io/badge/Go-1.25.6-00ADD8)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPL--3.0-green)](LICENSE)

> **Universal CDN IP Scanner for finding clean and low latency IPs**

This is an enhanced fork of [CFScanner](https://github.com/MortezaBashsiz/CFScanner/tree/main/golang) with significant improvements for performance, usability, and modern transport protocols. While originally designed for Cloudflare, it works effectively with any CDN to discover clean IPs with low latency.

## What Makes This Fork Special

**Latest Xray-Core v26.1.23** - Updated to the newest version with all security patches and performance improvements

**Modern TUI Interface** - Clean terminal interface with real-time statistics, worker status, and interactive controls

**XHTTP (SplitHTTP) Optimized** - Specifically optimized for the latest XHTTP transport protocol for better CDN compatibility

**Multi-threading Enhanced** - Proper multi-worker support with individual worker tracking and pause/resume functionality

**Universal CDN Support** - Works with Cloudflare, Fastly, KeyCDN, and other major CDN providers to find clean IPs

## Features

- **Real-time TUI Dashboard** with live statistics
- **Multi-threaded scanning** with individual worker status
- **Pause/Resume functionality** with P key
- **XHTTP (SplitHTTP) transport** for optimal CDN bypass
- **Cross-platform builds** (Linux, Windows, macOS)
- **Clean configuration management** with automatic cleanup
- **Comprehensive logging** with success/failure tracking
- **Progress tracking** with ETA calculations
- **Universal CDN compatibility** for finding low latency IPs

## Installation

### Quick Start

```bash
# Clone the repository
git clone https://github.com/ArchNets/CDN-IP-Scanner.git
cd CDN-IP-Scanner

# Build for your platform
go build -o cfscanner .
```

### Cross-Platform Builds

```bash
# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o cfscanner-windows-x64.exe .

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o cfscanner-macos-intel .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o cfscanner-macos-arm64 .
```

## Usage

### Basic Usage

```bash
# Scan CDN IP ranges with TUI interface
./cfscanner scan --subnets "104.16.0.0/12" --threads 8

# Use with custom config for any CDN
./cfscanner -c config.json --vpn -s ip_ranges.txt -t 16
```

### TUI Controls

- **P** - Pause/Resume scanning
- **Q / Esc** - Quit scanner
- **Ctrl+C** - Force quit

### Configuration

Create a `config.json` file:

```json
{
  "id": "your-uuid-here",
  "host": "your-cdn-server.example.com",
  "port": "443",
  "path": "/your-path",
  "serverName": "your-sni.example.com"
}
```

## TUI Interface

```
CDN IP SCANNER - CLEAN IP DISCOVERY
                    SCANNING  [████████████░░░░░░] 

Progress:     1,250 / 10,000  IPs  (12.5%)  |  Workers: 8/8
Success Rate: 15.2%  (Hits: 190  Fails: 1,060)
Speed:        45.3  IPs/sec  |  Elapsed: 27s  |  ETA: 3m 12s

WORKERS
Worker 1: 104.16.123.45    Worker 5: 172.67.89.12
Worker 2: 104.16.124.67    Worker 6: 172.67.90.34
Worker 3: 104.16.125.89    Worker 7: 172.67.91.56
Worker 4: 104.16.126.01    Worker 8: 172.67.92.78

RECENT ACTIVITY
14:23:45  |  104.16.123.45  |  Success
14:23:45  |  172.67.89.12   |  Failed
14:23:46  |  104.16.124.67  |  Success
```

## Command Line Options

| Option | Description | Example |
|--------|-------------|---------|
| `-c, --config` | Config file path | `-c config.json` |
| `-s, --subnets` | IP ranges file or CIDR | `-s ranges.txt` |
| `-t, --threads` | Number of workers | `-t 16` |
| `--vpn` | Enable VPN mode | `--vpn` |
| `--min-dl-speed` | Minimum download speed | `--min-dl-speed 100` |
| `--max-dl-latency` | Maximum latency (seconds) | `--max-dl-latency 2` |

## Build from Source

### Requirements

- Go 1.25.6 or later
- Git

### Dependencies

- **Xray-Core v26.1.23** - Latest version with XHTTP support
- **Bubble Tea** - Modern TUI framework
- **Lipgloss** - Terminal styling
- **Cobra** - CLI framework

### Build Process

```bash
# Clone repository
git clone https://github.com/ArchNets/CDN-IP-Scanner.git
cd CDN-IP-Scanner

# Install dependencies
go mod tidy

# Build
go build -o cfscanner .

# Run
./cfscanner scan --help
```

## Key Improvements Over Original

1. **Modern Transport Protocols**
   - XHTTP (SplitHTTP) as default transport
   - Better CDN compatibility and detection resistance

2. **Enhanced User Experience**
   - Clean TUI with real-time updates
   - Individual worker status tracking
   - Interactive pause/resume functionality

3. **Performance Optimizations**
   - Efficient multi-threading without race conditions
   - Optimized configuration file management
   - Reduced memory footprint

4. **Code Quality**
   - Clean, maintainable codebase
   - Proper error handling
   - Comprehensive logging

## CDN Compatibility

This scanner works effectively with major CDN providers:

- **Cloudflare** - Original target, fully tested
- **Fastly** - Compatible with XHTTP transport
- **KeyCDN** - Works with standard configurations
- **AWS CloudFront** - Supports clean IP discovery
- **Azure CDN** - Compatible with low latency scanning
- **Google Cloud CDN** - Works with proper configuration

## Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Credits

- **Original Project**: [CFScanner](https://github.com/MortezaBashsiz/CFScanner) by MortezaBashsiz
- **Enhanced Fork**: Developed and maintained by **Arch Net Team**
- **Xray-Core**: [XTLS Project](https://github.com/XTLS/Xray-core)

## Links

- [Original CFScanner](https://github.com/MortezaBashsiz/CFScanner/tree/main/golang)
- [Xray-Core Documentation](https://xtls.github.io/)
- [Transport Examples](vpn/TRANSPORT_EXAMPLES.md)

---

**Note**: This scanner is designed to find clean IPs with low latency across various CDN providers. Always ensure you comply with the terms of service of the services you're testing.