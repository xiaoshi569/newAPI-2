# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 配置 Go 代理（使用国内镜像）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制源代码和依赖配置
COPY go.mod ./
COPY . .

# 生成 go.sum 并下载依赖
RUN go mod tidy
RUN go mod download

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o router main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sync cmd/sync/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/router .
COPY --from=builder /app/sync .
COPY --from=builder /app/config.yaml .

# 暴露端口
EXPOSE 8080 9090

# 默认运行主程序
CMD ["./router"]
