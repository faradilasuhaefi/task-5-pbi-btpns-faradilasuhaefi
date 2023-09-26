package middlewares

import (
	"final-project-pbi-btpns/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			return
		}
		err := helpers.VerifyToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			return
		}
		context.Next()
	}
}
