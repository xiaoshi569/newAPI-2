package cache

import (
	"context"

	"api-router/internal/config"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端封装
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get 获取路由映射
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, "route:"+key).Result()
}

// Set 设置路由映射
func (r *RedisClient) Set(key, projectID string) error {
	return r.client.Set(r.ctx, "route:"+key, projectID, 0).Err()
}

// SetBatch 批量设置路由映射
func (r *RedisClient) SetBatch(keys []string, projectID string) error {
	pipe := r.client.Pipeline()

	for _, key := range keys {
		pipe.Set(r.ctx, "route:"+key, projectID, 0)
	}

	_, err := pipe.Exec(r.ctx)
	return err
}

// Delete 删除路由映射
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, "route:"+key).Err()
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Ping 测试连接
func (r *RedisClient) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Stats 获取统计信息
func (r *RedisClient) Stats() *redis.PoolStats {
	return r.client.PoolStats()
}
