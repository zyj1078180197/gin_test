package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "zyj.cn/docs"
	"zyj.cn/middlewares"
	"zyj.cn/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//路由规则
	r.GET("/ping", service.Ping)

	//公有方法
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
	//用户排行榜
	r.GET("/rank-list", service.GetRankList)

	//提交记录
	r.GET("/submit-list", service.GetSubmitList)

	// 分类列表
	r.GET("/category-list", service.GetCategoryList)

	//管理员私有方法
	authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	// 问题创建
	authAdmin.POST("/problem-create", service.ProblemCreate)
	// 问题修改
	authAdmin.PUT("/problem-modify", service.ProblemModify)
	// 分类创建
	authAdmin.POST("/category-create", service.CategoryCreate)
	// 分类修改
	authAdmin.PUT("/category-modify", service.CategoryModify)
	// 分类删除
	authAdmin.DELETE("/category-delete", service.CategoryDelete)
	// 用户私有方法
	authUser := r.Group("/user", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("/submit", service.Submit)
	return r

}
