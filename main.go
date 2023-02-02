package main

import (
	"github.com/gin-gonic/gin"
	"tic3001-go-server/config"
	"tic3001-go-server/router"
)

func main() {
	engine := gin.Default()
	router.Register(engine)
	_ = engine.Run(":" + config.Config.MustString("http.port", "8080"))
}
