package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omarelweshy/EcomMaster-user-service/internal/logger"
	"github.com/sirupsen/logrus"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		logger.GetLogger().WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"latency":  time.Since(startTime),
			"clientIP": c.ClientIP(),
		}).Info("request details")
	}
}
