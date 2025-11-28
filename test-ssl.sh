#!/bin/bash
# SSL 连接测试脚本

echo "=========================================="
echo "测试 MySQL SSL 连接"
echo "=========================================="

docker run -it --rm mysql:8.2 mysql \
  -h 38.175.200.23 \
  -P 3306 \
  -u new \
  -pnew \
  --ssl-mode=REQUIRED \
  -e "SHOW STATUS LIKE 'Ssl_cipher'; SHOW VARIABLES LIKE 'have_ssl';" \
  2>&1

echo ""
echo "=========================================="
echo "测试 PostgreSQL SSL 连接"
echo "=========================================="

docker run -it --rm postgres:latest psql \
  "host=8.148.201.143 port=5432 user=new dbname=new password=72Ddk4BfNzeYTCWp sslmode=require" \
  -c "SELECT version(); SHOW ssl; \conninfo" \
  2>&1

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="
echo "如果看到 SSL 相关信息，说明服务器支持 SSL"
echo "如果连接失败，说明服务器可能未启用 SSL"
