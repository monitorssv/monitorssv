package service

import (
	"github.com/gin-gonic/gin"
	ratelimiter "github.com/khaaleoo/gin-rate-limiter/core"
	"net/http"
	"slices"
	"time"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"https://localhost:3000",
	"https://monitorssv.xyz",
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if slices.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers, Content-Type")

		c.Header("X-Frame-Options", "DENY")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				monitorLog.Errorf("panic recover err: %v", err)

				ReturnErr(c, serverErrRes)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.URL.String()
		start := time.Now()
		c.Next()
		monitorLog.Infow("TraceLogger", "method", method, "cost", time.Since(start).String())
	}
}

func IpRageLimiter() gin.HandlerFunc {
	rateLimiterOption := ratelimiter.RateLimiterOption{
		Limit: 5,
		Burst: 200,
		Len:   10 * time.Minute,
	}

	rateLimiterMiddleware := ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip",
		Option:          rateLimiterOption,
	})
	return rateLimiterMiddleware
}
