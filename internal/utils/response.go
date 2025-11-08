package utils

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{"status": "success", "data": data})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"status": "error", "message": message})
}
