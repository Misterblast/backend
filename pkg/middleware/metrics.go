package middleware

import (
	"fmt"
	"time"

	metrics "github.com/ghulammuzz/misterblast/pkg/prom"
	"github.com/gofiber/fiber/v2"
)

func Metrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		route := c.Route()
		path := "unknown"

		if route != nil {
			path = route.Path
		}

		method := c.Method()
		status := c.Response().StatusCode()

		metrics.RequestCounter.WithLabelValues(path, method, fmt.Sprintf("%d", status)).Inc()

		metrics.RequestDuration.WithLabelValues(path, method).Observe(duration)

		if status >= 400 {
			metrics.ErrorCounter.WithLabelValues(path, method, fmt.Sprintf("%d", status)).Inc()
		}

		return err
	}
}
