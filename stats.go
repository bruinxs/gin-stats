package stats

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	ginLatencyMetric = "gin.latency"
	ginStatusMetric  = "gin.status"
	ginRequestMetric = "gin.request"
)

//Report from default metric registry
func Report() metrics.Registry {
	return metrics.DefaultRegistry
}

//RequestStats middleware
func RequestStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		req := metrics.GetOrRegisterMeter(ginRequestMetric, nil)
		req.Mark(1)

		since := time.Since(start)
		latency := metrics.GetOrRegisterTimer(ginLatencyMetric, nil)
		latency.Update(since)

		urlLatency := metrics.GetOrRegisterTimer(fmt.Sprintf("%s.%s", ginLatencyMetric, c.Request.URL.Path), nil)
		urlLatency.Update(since)

		status := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%d", ginStatusMetric, c.Writer.Status()), nil)
		status.Mark(1)
	}
}
