package middlewares

import (
	"GolangBase/task4/utils/token"
	"net/http"
	"github.com/gin-gonic/gin"
)

func JWtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.TokenValid(ctx)
		if err != nil {
			ctx.String(http.StatusUnauthorized,err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}