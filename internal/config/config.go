package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构
type Config struct {
	HTTP       HTTPConfig              `yaml:"http"`
	Redis      RedisConfig             `yaml:"redis"`
	LocalCache LocalCacheConfig        `yaml:"local_cache"`
	Projects   map[string]ProjectConfig `yaml:"projects"`
	Sync       SyncConfig              `yaml:"sync"`
	Metrics    MetricsConfig           `yaml:"metrics"`
	Log        LogConfig               `yaml:"log"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// LocalCacheConfig 本地缓存配置
type LocalCacheConfig struct {
	MaxSize int           `yaml:"max_size"`
	TTL     time.Duration `yaml:"ttl"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Backends []string       `yaml:"backends"`
	Database DatabaseConfig `yaml:"database"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string `yaml:"type"`          // 数据库类型: mysql 或 postgres
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	SSLMode      string `yaml:"ssl_mode"`      // SSL 模式: disable, require, verify-ca, verify-full (PostgreSQL) 或 true/false/skip-verify (MySQL)
}

// SyncConfig 同步任务配置
type SyncConfig struct {
	Interval  time.Duration `yaml:"interval"`
	BatchSize int           `yaml:"batch_size"`
	Enabled   bool          `yaml:"enabled"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DSN 生成数据库连接字符串
func (d *DatabaseConfig) DSN() string {
	port := fmt.Sprintf("%d", d.Port)

	switch d.Type {
	case "postgres", "postgresql":
		// PostgreSQL: host=localhost port=5432 user=xxx password=xxx dbname=xxx sslmode=xxx
		sslMode := d.SSLMode
		if sslMode == "" {
			sslMode = "disable" // 默认不启用 SSL
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			d.Host, port, d.User, d.Password, d.DBName, sslMode)
	case "mysql", "":
		// MySQL: user:password@tcp(host:port)/dbname?parseTime=true&tls=xxx
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			d.User, d.Password, d.Host, port, d.DBName)

		// 添加 SSL 配置
		if d.SSLMode != "" && d.SSLMode != "disable" && d.SSLMode != "false" {
			if d.SSLMode == "true" || d.SSLMode == "require" {
				dsn += "&tls=true"
			} else if d.SSLMode == "skip-verify" {
				dsn += "&tls=skip-verify"
			}
			// 注意：自定义 TLS 配置名称需要先通过 mysql.RegisterTLSConfig() 注册
			// 为避免连接失败，这里只支持内置的 true/skip-verify 两种模式
		}
		return dsn
	default:
		// 默认使用MySQL格式
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			d.User, d.Password, d.Host, port, d.DBName)
	}
}

// Driver 返回数据库驱动名称
func (d *DatabaseConfig) Driver() string {
	switch d.Type {
	case "postgres", "postgresql":
		return "postgres"
	case "mysql", "":
		return "mysql"
	default:
		return "mysql"
	}
}
