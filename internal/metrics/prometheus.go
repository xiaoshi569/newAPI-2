package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 请求计数
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "router_requests_total",
			Help: "Total number of requests",
		},
		[]string{"project", "result", "status"},
	)

	// 请求延迟
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "router_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"project"},
	)

	// 路由查询延迟
	lookupDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "router_lookup_duration_seconds",
			Help:    "Lookup duration in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
		},
		[]string{"result"},
	)

	// 缓存命中
	cacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "router_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"level"}, // local, redis, database
	)

	// 缓存未命中
	cacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "router_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"level"},
	)
)

func init() {
	// 注册指标
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(lookupDuration)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
}

// Handler 返回Prometheus HTTP处理器
func Handler() http.Handler {
	return promhttp.Handler()
}

// RecordRequest 记录请求
func RecordRequest(project, result string, status int) {
	requestsTotal.WithLabelValues(project, result, strconv.Itoa(status)).Inc()
}

// RecordRequestStart 记录请求开始时间
func RecordRequestStart() time.Time {
	return time.Now()
}

// RecordRequestDuration 记录请求延迟
func RecordRequestDuration(start time.Time, project string) {
	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(project).Observe(duration)
}

// RecordLookupDuration 记录查询延迟
func RecordLookupDuration(start time.Time, result string) {
	duration := time.Since(start).Seconds()
	lookupDuration.WithLabelValues(result).Observe(duration)
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(level string) {
	cacheHits.WithLabelValues(level).Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(level string) {
	cacheMisses.WithLabelValues(level).Inc()
}
