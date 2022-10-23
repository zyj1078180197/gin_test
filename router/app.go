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

	//问题列表
	r.GET("/problem-list", service.GetProblemList)
	//问题详情
	r.GET("/problem-detail", service.GetProblemDetail)

	//用户登录
	r.POST("/login", service.Login)
	//用户注册
	r.POST("/register", service.Register)
    //用户详情
    r.GET("/user-detail", service.GetUserDetail)
	//发送验证码
	r.POST("/send-code", service.SendCode)

	//提交记录
	r.GET("/submit-list", service.GetSubmitList)


	return r

}
