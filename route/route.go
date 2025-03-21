package route

import (
	"cf-analytics-server/config"
	"cf-analytics-server/controller"
	"cf-analytics-server/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func Start() {
	mode := gin.ReleaseMode
	if config.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	r := gin.Default()

	r.Use(middleware.Cors())

	r.POST("/statistics", controller.Statistics)

	err := r.Run("0.0.0.0:4000")
	if err != nil {
		log.Fatalln("服务器启动失败！错误原因：" + err.Error())
	}
}
