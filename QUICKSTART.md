# 快速开始指南

## 🚀 5分钟部署

### 第1步：配置文件

编辑 `config.yaml`，填入你的项目信息：

```yaml
projects:
  project_a:
    backends:
      - "https://your-project-a.com"    # 项目A的地址
    database:
      type: "mysql"                      # 数据库类型: mysql 或 postgres
      host: "localhost"
      port: 3306                         # MySQL: 3306, PostgreSQL: 5432
      user: "readonly"                   # 只读用户即可
      password: "your_password"
      dbname: "one_api"

  project_b:
    backends:
      - "https://your-project-b.com"    # 项目B的地址
    database:
      type: "mysql"                      # 可以与项目A使用不同的数据库类型
      host: "localhost"
      port: 3306
      user: "readonly"
      password: "your_password"
      dbname: "one_api"
```

### 第2步：启动服务

```bash
# 使用 Docker Compose 一键启动（推荐）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 等待服务启动（约10秒）
```

### 第3步：验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 查看监控指标
curl http://localhost:9090/metrics

# 测试路由（使用你的实际API Key）
curl -H "Authorization: Bearer sk-your-key-here" \
     http://localhost:8080/v1/models
```

### 第4步：配置域名（可选）

使用 Nginx 作为反向代理：

```nginx
# /etc/nginx/sites-available/api-router
server {
    listen 80;
    server_name api.your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时配置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
}
```

---

## 📊 验证数据同步

```bash
# 进入 Redis 查看映射数据
docker exec -it api-router-redis redis-cli

# 查看某个key的映射
127.0.0.1:6379> GET route:your-key-here

# 查看所有路由映射的数量
127.0.0.1:6379> KEYS route:*
```

---

## 🔧 常见问题

### 1. Redis 连接失败

**错误**: `Failed to connect to Redis`

**解决**: 确保 Redis 已启动且端口正确
```bash
docker-compose ps redis
```

### 2. 数据库连接失败

**错误**: `Failed to connect to database`

**解决**:
- 检查数据库地址、端口、用户名、密码
- 检查数据库类型配置（`type: mysql` 或 `type: postgres`）
- 确保数据库允许远程连接
- 测试连接：
  ```bash
  # MySQL
  mysql -h host -u user -p

  # PostgreSQL
  psql -h host -U user -d dbname
  ```

**支持的数据库类型**:
- ✅ MySQL 5.7+, 8.0+（端口默认 3306）
- ✅ PostgreSQL 12+, 13+, 14+, 15+（端口默认 5432）
- ✅ 每个项目可以使用不同的数据库类型

### 3. Key 查询失败

**错误**: `invalid API key`

**原因**:
- Key 还未同步到 Redis（等待 0-5 分钟）
- Key 在数据库中不存在或已删除
- Key 状态不是启用状态（status != 1）

**解决**:
```bash
# 手动触发同步（重启同步容器）
docker-compose restart sync

# 查看同步日志
docker-compose logs sync
```

### 4. 转发失败

**错误**: `backend unavailable`

**原因**:
- 后端项目不可访问
- 网络不通

**解决**:
- 检查 backend URL 是否正确
- 测试后端连接：`curl https://your-project.com`

---

## 📈 性能优化

### 1. 增加实例数量

编辑 `docker-compose.yml`：

```yaml
router:
  deploy:
    replicas: 3  # 增加到3个实例
```

### 2. 调整缓存大小

编辑 `config.yaml`：

```yaml
local_cache:
  max_size: 200000  # 增加到20万
  ttl: 600s         # 延长到10分钟
```

### 3. 调整同步频率

```yaml
sync:
  interval: 3m      # 改为3分钟同步一次
```

---

## 🔄 更新部署

```bash
# 停止服务
docker-compose down

# 拉取最新代码 / 修改配置

# 重新构建并启动
docker-compose up -d --build

# 查看日志确认启动成功
docker-compose logs -f router
```

---

## 📊 监控告警

### Grafana Dashboard

导入以下 Prometheus 查询：

**QPS（每秒请求数）：**
```promql
rate(router_requests_total[1m])
```

**P99 延迟：**
```promql
histogram_quantile(0.99, rate(router_lookup_duration_seconds_bucket[5m]))
```

**缓存命中率：**
```promql
sum(rate(router_cache_hits_total[5m])) /
sum(rate(router_cache_hits_total[5m]) + rate(router_cache_misses_total[5m]))
```

---

## ❓ 需要帮助？

1. 查看日志：`docker-compose logs -f`
2. 检查服务状态：`docker-compose ps`
3. 进入容器调试：`docker exec -it api-router sh`

**生产环境建议：**
- 使用独立的 Redis Cluster（而非单机 Redis）
- 配置健康检查和自动重启
- 部署多个中间件实例（3+）
- 配置负载均衡（Nginx/HAProxy）
- 设置 Prometheus + Grafana 监控
- 配置告警规则（AlertManager）
