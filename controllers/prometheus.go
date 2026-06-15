package controllers

import (
	"bluebell/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method

		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		metrics.ObserveHTTPRequest(method, path, status, duration)
	}
}
