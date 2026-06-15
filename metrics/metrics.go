package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

const namespace = "bluebell"

var (
	cpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "runtime",
		Name:      "cpu_usage_percent",
		Help:      "Current CPU usage percentage.",
	}, []string{"instance"})

	memoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "runtime",
		Name:      "memory_usage_percent",
		Help:      "Current memory usage percentage.",
	}, []string{"instance"})

	goroutineNum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "runtime",
		Name:      "goroutine_num",
		Help:      "Current goroutine count.",
	}, []string{"instance"})

	processNum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "runtime",
		Name:      "process_num",
		Help:      "Current process count on host.",
	}, []string{"instance"})

	httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "request_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})
)

func Init(instance string) {
	if instance == "" {
		instance = "local"
	}
	prometheus.MustRegister(
		cpuUsage,
		memoryUsage,
		goroutineNum,
		processNum,
		httpRequestTotal,
		httpRequestDuration,
	)
	startRuntimeCollector(instance)
}

func ObserveHTTPRequest(method, path, status string, durationSeconds float64) {
	httpRequestTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path, status).Observe(durationSeconds)
}

func startRuntimeCollector(instance string) {
	go func() {
		collectRuntimeMetrics(instance)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			collectRuntimeMetrics(instance)
		}
	}()
}

func collectRuntimeMetrics(instance string) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		zap.L().Error("get cpu usage failed", zap.Error(err))
		cpuUsage.WithLabelValues(instance).Set(0)
	} else {
		cpuUsage.WithLabelValues(instance).Set(cpuPercent[0])
	}

	memoryPercent, err := mem.VirtualMemory()
	if err != nil {
		zap.L().Error("get memory usage failed", zap.Error(err))
		memoryUsage.WithLabelValues(instance).Set(0)
	} else {
		memoryUsage.WithLabelValues(instance).Set(memoryPercent.UsedPercent)
	}

	goroutineNum.WithLabelValues(instance).Set(float64(runtime.NumGoroutine()))

	processes, err := process.Processes()
	if err != nil {
		zap.L().Error("get process count failed", zap.Error(err))
		processNum.WithLabelValues(instance).Set(0)
	} else {
		processNum.WithLabelValues(instance).Set(float64(len(processes)))
	}
}
