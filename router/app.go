package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "zyj.cn/docs"
	"zyj.cn/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//路由规则
	r.GET("/ping", service.Ping)

	r.GET("/problem-list", service.GetProblemList)

	return r

}
