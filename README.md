# API Router Middleware - 独立路由中间件

## 📋 项目说明

这是一个**完全独立**的 API 路由中间件，用于在多个 API 项目之间进行智能路由转发。

**核心功能：**
- ✅ 根据 API Key 自动路由到不同的后端项目
- ✅ 三层缓存架构（本地内存 + Redis + 数据库）
- ✅ 高性能（10,000+ QPS/实例，< 2ms 延迟）
- ✅ 零侵入（不修改任何现有项目代码）

---

## 🏗️ 项目结构

```
api-router-middleware/
├── README.md                 # 项目说明
├── go.mod                    # Go模块定义
├── go.sum                    # 依赖版本锁定
├── main.go                   # 主程序入口
├── config.yaml               # 配置文件
├── Dockerfile                # Docker镜像构建
├── docker-compose.yml        # 容器编排
│
├── cmd/
│   └── sync/
│       └── main.go          # 数据同步工具
│
├── internal/
│   ├── config/
│   │   └── config.go        # 配置加载
│   ├── router/
│   │   ├── handler.go       # 请求处理
│   │   ├── extractor.go     # Key提取器
│   │   └── proxy.go         # 反向代理
│   ├── cache/
│   │   ├── redis.go         # Redis缓存
│   │   └── local.go         # 本地缓存
│   ├── database/
│   │   └── client.go        # 数据库客户端
│   └── metrics/
│       └── prometheus.go    # 监控指标
│
└── scripts/
    └── deploy.sh            # 部署脚本
```

---

## 🚀 快速开始

### 1. 配置文件

编辑 `config.yaml`：

```yaml
redis:
  addr: "localhost:6379"
  password: ""
  db: 0

projects:
  project_a:
    backends:
      - "https://your-project-a.com"
    database:
      type: "mysql"              # 数据库类型: mysql 或 postgres
      host: "db-a.example.com"
      port: 3306                 # MySQL: 3306, PostgreSQL: 5432
      user: "readonly"
      password: "your_password"
      dbname: "one_api"

  project_b:
    backends:
      - "https://your-project-b.com"
    database:
      type: "postgres"           # 支持混用不同数据库类型
      host: "db-b.example.com"
      port: 5432
      user: "readonly"
      password: "your_password"
      dbname: "one_api"
```

### 2. Docker Compose 部署（推荐）

```bash
# 启动所有服务（Redis + 中间件 + 同步任务）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 3. 手动部署

```bash
# 安装依赖
go mod download

# 运行数据同步（首次）
go run cmd/sync/main.go

# 启动中间件
go run main.go
```

---

## 📊 性能指标

- **QPS**: 10,000+ /实例
- **延迟**: P99 < 2ms（中间件部分）
- **资源**: 4核8G，800MB内存

---

## 🔧 监控

访问 `http://localhost:9090/metrics` 查看 Prometheus 指标

---

## 💾 数据库支持

### 支持的数据库类型

- ✅ **MySQL** (5.7+, 8.0+)
- ✅ **PostgreSQL** (12+, 13+, 14+, 15+)

### 混合数据库配置

**每个项目可以使用不同的数据库类型**，例如：

```yaml
projects:
  project_a:
    database:
      type: "mysql"        # 项目A使用MySQL
      port: 3306
      # ...

  project_b:
    database:
      type: "postgres"     # 项目B使用PostgreSQL
      port: 5432
      # ...
```

### 数据库连接格式

**MySQL:**
```
user:password@tcp(host:port)/dbname?parseTime=true
```

**PostgreSQL:**
```
host=xxx port=5432 user=xxx password=xxx dbname=xxx sslmode=disable
```

### 数据库权限要求

中间件**只需要只读权限**：

```sql
-- MySQL 授权示例
GRANT SELECT ON one_api.tokens TO 'readonly'@'%';

-- PostgreSQL 授权示例
GRANT SELECT ON TABLE tokens TO readonly;
```

---

## 📝 工作原理

```
客户端请求
    ↓
提取 API Key
    ↓
三层缓存查询:
  L1: 本地内存 (0.1ms)  ← 90% 命中
  L2: Redis (0.5ms)     ← 9.9% 命中
  L3: 数据库 (10ms)     ← 0.1% 命中
    ↓
反向代理转发到对应项目
    ↓
项目自己验证并处理
```

---

## 📞 联系方式

有问题请联系项目维护者
