# Zeo++ Analysis API - Go Implementation

A high-performance, production-ready Go service that provides RESTful API endpoints for Zeo++ structure analysis.

## Features

- **High Concurrency**: Worker pool pattern with rate limiting
- **Smart Caching**: File-based caching with TTL and sharding
- **Production Ready**: Docker support, health checks, graceful shutdown
- **Comprehensive**: All 9 Zeo++ analysis endpoints
- **Type Safe**: Full Go type system with validation
- **Performance Optimized**: Streaming file uploads, buffer pools

## Quick Start

### Using Docker (Recommended)

```bash
# Clone and build
git clone <repository>
cd zeo++_api_go

# Start with Docker Compose
docker-compose up --build

# Service will be available at http://localhost:8080
```

### Local Development

```bash
# Install dependencies
go mod download

# Install Zeo++
# On macOS: brew install zeo++
# On Linux: follow instructions from http://www.zeoplusplus.org/

# Run server
go run cmd/server/main.go
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/pore_diameter` | POST | Calculate pore diameters |
| `/api/surface_area` | POST | Calculate accessible surface area |
| `/api/accessible_volume` | POST | Calculate accessible volume |
| `/api/probe_volume` | POST | Calculate probe-occupiable volume |
| `/api/channel_analysis` | POST | Analyze channel dimensionality |
| `/api/framework_info` | POST | Get framework information |
| `/api/pore_size_dist/download` | POST | Download pore size distribution |
| `/api/blocking_spheres` | POST | Generate blocking spheres |
| `/api/open_metal_sites` | POST | Count open metal sites |
| `/health` | GET | Health check |

## Usage Examples

### Calculate Pore Diameter

```bash
curl -X POST http://localhost:8080/api/pore_diameter \
  -F "structure_file=@/path/to/structure.cif" \
  -F "ha=true"
```

### Calculate Surface Area

```bash
curl -X POST http://localhost:8080/api/surface_area \
  -F "structure_file=@/path/to/structure.cif" \
  -F "probe_radius=1.21" \
  -F "samples=2000" \
  -F "ha=true"
```

### Download Pore Size Distribution

```bash
curl -X POST http://localhost:8080/api/pore_size_dist/download \
  -F "structure_file=@/path/to/structure.cif" \
  -F "probe_radius=1.21" \
  -F "samples=50000" \
  -o pore_size_distribution.psd
```

## Configuration

### Environment Variables

- `ZEO_EXEC_PATH`: Path to Zeo++ executable (default: "network")
- `ZEO_WORKSPACE`: Working directory for temporary files (default: "./workspace")
- `ENABLE_CACHE`: Enable/disable caching (default: true)

### Configuration File

Edit `config/config.yaml`:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  read_timeout: 30s
  write_timeout: 30s

zeo:
  executable_path: "network"
  workdir: "./workspace"
  timeout: 5m

concurrency:
  max_workers: 0           # 0 = runtime.NumCPU()
  max_queue_size: 1000
  rate_limit_per_ip: 10    # requests per second
  max_file_size: 100MB

cache:
  enabled: true
  ttl: 3600s              # 1 hour
  shards: 32
```

## Supported File Formats

- `.cif` - Crystallographic Information File
- `.cssr` - CSSR format
- `.v1` - V1 format
- `.arc` - ARC format
- `.cif.gz` - Gzipped CIF files
- `.cssr.gz` - Gzipped CSSR files

## Performance

- **Concurrency**: Configurable worker pool (default: CPU cores)
- **Caching**: SHA256-based caching with TTL
- **Rate Limiting**: Per-IP and global limits
- **Memory**: Buffer pools for efficient file handling
- **Monitoring**: Health checks and metrics endpoints

## Development

### Project Structure

```bash
zeo++_api_go/
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and middleware
│   ├── core/           # Core business logic
│   │   ├── runner/     # Zeo++ execution
│   │   ├── cache/      # Caching system
│   │   ├── pool/       # Worker pool
│   │   └── parser/     # Output parsers
│   └── utils/          # Utilities
├── config/             # Configuration files
├── tests/              # Test files
├── Dockerfile          # Container configuration
└── docker-compose.yml  # Docker Compose setup
```

### Testing

```bash
# Run tests
go test ./...

# Run with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

### Load Testing

```bash
# Install k6
brew install k6

# Run load test
k6 run tests/load_test.js
```

## Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### Cache Statistics

Cache statistics are logged periodically. Check logs for:
- Cache hit/miss ratios
- Worker pool utilization
- Request rates and latencies

## Troubleshooting

### Common Issues

1. **Zeo++ not found**: Ensure Zeo++ is installed and in PATH
2. **Permission denied**: Check workspace directory permissions
3. **Memory issues**: Adjust max_file_size and max_workers
4. **Timeouts**: Increase zeo.timeout in config

### Logs

Logs are output in JSON format. Set `LOG_LEVEL=debug` for verbose logging.

### Debug Mode

```bash
# Enable debug mode
export GIN_MODE=debug
go run cmd/server/main.go
```

## License

MIT License - see LICENSE file for details.

---

[中文版](README-cn.md)