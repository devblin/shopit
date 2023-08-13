package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorsConfig struct {
	IsDevMode      bool
	AllowedOrigins []string
}

func Cors(config CorsConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if config.IsDevMode || ctx.Request.Method == http.MethodOptions {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}

		if !config.IsDevMode {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
