# 项目文件清单

## 📁 项目结构

```
api-router-middleware/
├── 📄 README.md                    # 项目说明
├── 📄 QUICKSTART.md                # 快速开始指南
├── 📄 FILES.md                     # 本文件（文件清单）
├── 📄 go.mod                       # Go模块定义
├── 📄 go.sum                       # 依赖锁定（运行后自动生成）
├── 📄 config.yaml                  # 配置文件（需手动创建）
├── 📄 config.yaml.example          # 配置示例
├── 📄 .gitignore                   # Git忽略文件
├── 📄 Dockerfile                   # Docker镜像构建
├── 📄 docker-compose.yml           # 容器编排
│
├── 📄 main.go                      # 主程序入口
│
├── cmd/
│   └── sync/
│       └── 📄 main.go             # 数据同步工具主程序
│
└── internal/
    ├── cache/
    │   ├── 📄 redis.go            # Redis缓存客户端
    │   └── 📄 local.go            # 本地内存缓存
    │
    ├── config/
    │   └── 📄 config.go           # 配置加载和解析
    │
    ├── router/
    │   ├── 📄 handler.go          # 路由处理核心逻辑
    │   ├── 📄 extractor.go        # API Key提取器
    │   └── 📄 sync.go             # 数据同步任务
    │
    └── metrics/
        └── 📄 prometheus.go       # Prometheus监控指标
```

---

## 📄 文件说明

### 📚 文档文件

| 文件 | 说明 |
|------|------|
| `README.md` | 项目整体介绍、功能说明、架构图 |
| `QUICKSTART.md` | 快速部署指南、配置说明、常见问题 |
| `FILES.md` | 本文件，完整的文件清单和说明 |

### ⚙️ 配置文件

| 文件 | 说明 |
|------|------|
| `go.mod` | Go模块定义，包含项目依赖 |
| `go.sum` | 依赖版本锁定（运行`go mod download`后生成） |
| `config.yaml` | **运行时配置文件**（需要手动创建，参考example） |
| `config.yaml.example` | 配置文件模板，包含所有配置项说明 |
| `.gitignore` | Git忽略规则（不提交敏感配置） |

### 🐳 部署文件

| 文件 | 说明 |
|------|------|
| `Dockerfile` | Docker镜像构建文件（多阶段构建） |
| `docker-compose.yml` | 容器编排文件（Redis + 中间件 + 同步任务） |

### 💻 Go 代码文件

#### 主程序
- `main.go` - 中间件主程序入口，启动HTTP服务和监控
- `cmd/sync/main.go` - 数据同步工具，定时从数据库同步Key到Redis

#### 缓存模块 (internal/cache/)
- `redis.go` - Redis客户端封装，提供Get/Set/Batch操作
- `local.go` - 本地内存缓存（LRU），带自动过期清理

#### 配置模块 (internal/config/)
- `config.go` - 配置文件加载和结构定义

#### 路由模块 (internal/router/)
- `handler.go` - **核心**：三层缓存查询 + 反向代理转发
- `extractor.go` - 从请求中提取API Key（支持多种格式）
- `sync.go` - 数据同步逻辑（从数据库读取Key，批量写入Redis）

#### 监控模块 (internal/metrics/)
- `prometheus.go` - Prometheus指标定义和记录

---

## 🔑 关键文件详解

### 1. `main.go` - 主程序

**功能：**
- 加载配置文件
- 初始化Redis和本地缓存
- 创建路由处理器
- 启动HTTP服务（端口8080）
- 启动监控服务（端口9090）
- 启动数据同步任务
- 优雅关闭

**核心代码：**
```go
r.Use(routeHandler.RouteMiddleware())  // 所有请求走路由中间件
```

### 2. `internal/router/handler.go` - 路由核心

**功能：**
- 提取API Key
- 三层缓存查询（本地 → Redis → 数据库）
- 反向代理转发
- 监控指标记录

**查询逻辑：**
```
L1: 本地缓存 (0.1ms) → 90% 命中
  ↓ 未命中
L2: Redis (0.5ms) → 9.9% 命中
  ↓ 未命中
L3: 数据库 (10ms) → 0.1% 命中
```

### 3. `internal/router/sync.go` - 数据同步

**功能：**
- 定时任务（默认5分钟）
- 连接各项目数据库（只读）
- 查询有效的Token
- 批量写入Redis
- 日志记录

**SQL查询：**
```sql
SELECT `key`
FROM tokens
WHERE deleted_at IS NULL
  AND status = 1
```

### 4. `docker-compose.yml` - 容器编排

**包含3个服务：**
1. `redis` - 独立Redis实例（存储路由映射）
2. `router` - 中间件HTTP服务
3. `sync` - 数据同步任务

**一键启动所有服务：**
```bash
docker-compose up -d
```

---

## 💾 数据库支持

### 支持的数据库类型

中间件支持以下数据库类型，**每个项目可以使用不同的数据库**：

| 数据库类型 | 版本要求 | 默认端口 | 配置值 |
|-----------|---------|---------|--------|
| MySQL | 5.7+, 8.0+ | 3306 | `type: "mysql"` |
| PostgreSQL | 12+, 13+, 14+, 15+ | 5432 | `type: "postgres"` |

### 配置示例

```yaml
projects:
  project_a:
    database:
      type: "mysql"              # 项目A使用MySQL
      host: "db-a.example.com"
      port: 3306
      user: "readonly"
      password: "xxx"
      dbname: "one_api"

  project_b:
    database:
      type: "postgres"           # 项目B使用PostgreSQL
      host: "db-b.example.com"
      port: 5432
      user: "readonly"
      password: "xxx"
      dbname: "one_api"
```

### 数据库权限

中间件**只需要只读权限**，只查询 `tokens` 表：

```sql
-- MySQL 授权
GRANT SELECT ON one_api.tokens TO 'readonly'@'%';

-- PostgreSQL 授权
GRANT SELECT ON TABLE tokens TO readonly;
```

### 依赖的驱动

- MySQL: `github.com/go-sql-driver/mysql`
- PostgreSQL: `github.com/lib/pq`

---

## 📝 使用步骤

### 第1步：创建配置文件

```bash
cp config.yaml.example config.yaml
vim config.yaml  # 填入实际配置
```

### 第2步：启动服务

```bash
docker-compose up -d
```

### 第3步：验证

```bash
# 健康检查
curl http://localhost:8080/health

# 测试路由
curl -H "Authorization: Bearer sk-your-key" \
     http://localhost:8080/v1/models
```

---

## 🔧 开发说明

### 本地开发

```bash
# 安装依赖
go mod download

# 运行主程序
go run main.go

# 运行同步任务
go run cmd/sync/main.go
```

### 构建二进制

```bash
# 构建中间件
go build -o router main.go

# 构建同步工具
go build -o sync cmd/sync/main.go
```

### 添加新项目

1. 编辑 `config.yaml`，添加新项目配置
2. 重启服务：`docker-compose restart`

---

## 📊 监控指标

访问 `http://localhost:9090/metrics` 查看所有指标：

- `router_requests_total` - 请求总数
- `router_lookup_duration_seconds` - 查询延迟
- `router_cache_hits_total` - 缓存命中数
- `router_cache_misses_total` - 缓存未命中数

---

## 🆘 故障排查

### 查看日志
```bash
docker-compose logs -f router
docker-compose logs -f sync
docker-compose logs -f redis
```

### 进入容器
```bash
docker exec -it api-router sh
docker exec -it api-router-redis redis-cli
```

### 重启服务
```bash
docker-compose restart router
docker-compose restart sync
```

---

**项目维护者**: [填写你的名字]
**最后更新**: 2025-01-20
