# Zeo++ 分析 API - Go 实现

一个高性能、生产就绪的 Go 服务，提供用于 Zeo++ 结构分析的 RESTful API 端点。

## 功能

- **高并发**: 使用工作池模式和速率限制
- **智能缓存**: 基于文件的缓存，支持 TTL 和分片
- **生产就绪**: 支持 Docker、健康检查、优雅关闭
- **全面**: 提供所有 9 个 Zeo++ 分析端点
- **类型安全**: 完整的 Go 类型系统和验证
- **性能优化**: 支持流式文件上传和缓冲池

## 快速开始

### 使用 Docker（推荐）

```bash
# 克隆并构建
git clone <repository>
cd zeo++_api_go

# 使用 Docker Compose 启动
docker-compose up --build

# 服务将可通过 http://localhost:8080 访问
```

### 本地开发

```bash
# 安装依赖
go mod download

# 安装 Zeo++
# 在 macOS 上: brew install zeo++
# 在 Linux 上: 按照 http://www.zeoplusplus.org/ 的说明操作

# 运行服务器
go run cmd/server/main.go
```

## API 端点

| 端点 | 方法 | 描述 |
|----------|--------|-------------|
| `/api/pore_diameter` | POST | 计算孔径 |
| `/api/surface_area` | POST | 计算可达表面积 |
| `/api/accessible_volume` | POST | 计算可达体积 |
| `/api/probe_volume` | POST | 计算探针可占体积 |
| `/api/channel_analysis` | POST | 分析通道维度 |
| `/api/framework_info` | POST | 获取框架信息 |
| `/api/pore_size_dist/download` | POST | 下载孔径分布 |
| `/api/blocking_spheres` | POST | 生成阻塞球 |
| `/api/open_metal_sites` | POST | 统计开放金属位点 |
| `/health` | GET | 健康检查 |

## 使用示例

### 计算孔径

```bash
curl -X POST http://localhost:8080/api/pore_diameter \
  -F "structure_file=@/path/to/structure.cif" \
  -F "ha=true"
```

### 计算表面积

```bash
curl -X POST http://localhost:8080/api/surface_area \
  -F "structure_file=@/path/to/structure.cif" \
  -F "probe_radius=1.21" \
  -F "samples=2000" \
  -F "ha=true"
```

### 下载孔径分布

```bash
curl -X POST http://localhost:8080/api/pore_size_dist/download \
  -F "structure_file=@/path/to/structure.cif" \
  -F "probe_radius=1.21" \
  -F "samples=50000" \
  -o pore_size_distribution.psd
```

## 配置

### 环境变量

- `ZEO_EXEC_PATH`: Zeo++ 可执行文件路径（默认: "network"）
- `ZEO_WORKSPACE`: 用于临时文件的工作目录（默认: "./workspace"）
- `ENABLE_CACHE`: 启用/禁用缓存（默认: true）

### 配置文件

编辑 `config/config.yaml`:

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
  rate_limit_per_ip: 10    # 每秒请求数
  max_file_size: 100MB

cache:
  enabled: true
  ttl: 3600s              # 1 小时
  shards: 32
```

## 支持的文件格式

- `.cif` - 晶体学信息文件
- `.cssr` - CSSR 格式
- `.v1` - V1 格式
- `.arc` - ARC 格式
- `.cif.gz` - 压缩的 CIF 文件
- `.cssr.gz` - 压缩的 CSSR 文件

## 性能

- **并发**: 可配置的工作池（默认: CPU 核心数）
- **缓存**: 基于 SHA256 的缓存，支持 TTL
- **速率限制**: 按 IP 和全局限制
- **内存**: 使用缓冲池高效处理文件
- **监控**: 健康检查和指标端点

## 开发

### 项目结构

```bash
zeo++_api_go/
├── cmd/server/          # 应用程序入口点
├── internal/
│   ├── api/            # HTTP 处理程序和中间件
│   ├── core/           # 核心业务逻辑
│   │   ├── runner/     # Zeo++ 执行
│   │   ├── cache/      # 缓存系统
│   │   ├── pool/       # 工作池
│   │   └── parser/     # 输出解析器
│   └── utils/          # 工具
├── config/             # 配置文件
├── tests/              # 测试文件
├── Dockerfile          # 容器配置
└── docker-compose.yml  # Docker Compose 设置
```

### 测试

```bash
# 运行测试
go test ./...

# 启用竞争检测运行
go test -race ./...

# 基准测试
go test -bench=. ./...
```

### 压力测试

```bash
# 安装 k6
brew install k6

# 运行压力测试
k6 run tests/load_test.js
```

## 监控

### 健康检查

```bash
curl http://localhost:8080/health
```

### 缓存统计

缓存统计会定期记录到日志中。检查日志以获取:

- 缓存命中/未命中比率
- 工作池利用率
- 请求速率和延迟

## 故障排除

### 常见问题

1. **找不到 Zeo++**: 确保 Zeo++ 已安装并在 PATH 中
2. **权限被拒绝**: 检查工作目录权限
3. **内存问题**: 调整 max_file_size 和 max_workers
4. **超时**: 增加配置中的 zeo.timeout

### 日志

日志以 JSON 格式输出。设置 `LOG_LEVEL=debug` 以启用详细日志。

### 调试模式

```bash
# 启用调试模式
export GIN_MODE=debug
go run cmd/server/main.go
```

## 许可证

MIT 许可证 - 详情见 LICENSE 文件。

---

[English Version](README.md)
