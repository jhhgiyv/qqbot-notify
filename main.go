package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jhhgiyv/qqbot-notify/config"
)

func main() {
	route := gin.Default()
	route.POST("/", notify)
}

func notify(context *gin.Context) {
}
