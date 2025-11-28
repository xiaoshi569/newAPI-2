package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-router/internal/cache"
	"api-router/internal/config"
	"api-router/internal/metrics"
	"api-router/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化 Redis
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// 3. 初始化本地缓存
	localCache := cache.NewLocalCache(cfg.LocalCache.MaxSize, cfg.LocalCache.TTL)

	// 4. 初始化路由处理器
	routeHandler, err := router.NewHandler(cfg, redisClient, localCache)
	if err != nil {
		log.Fatalf("Failed to create router handler: %v", err)
	}

	// 5. 启动监控服务
	if cfg.Metrics.Enabled {
		go startMetricsServer(cfg.Metrics.Port, cfg.Metrics.Path)
	}

	// 6. 启动数据同步任务（如果启用）
	if cfg.Sync.Enabled {
		go router.StartSyncTask(cfg, redisClient)
	}

	// 7. 创建 Gin 服务器
	if cfg.Log.Level == "info" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// 所有其他请求都走路由中间件
	r.Use(routeHandler.RouteMiddleware())

	// 8. 启动 HTTP 服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// 优雅关闭
	go func() {
		log.Printf("Starting API Router on port %d", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// startMetricsServer 启动 Prometheus 监控服务
func startMetricsServer(port int, path string) {
	mux := http.NewServeMux()
	mux.Handle(path, metrics.Handler())

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting metrics server on %s%s", addr, path)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("Failed to start metrics server: %v", err)
	}
}
