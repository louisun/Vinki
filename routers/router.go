package routers

import (
	"github.com/vinki/routers/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// api 接口
	v1 := r.Group("/api/v1")

	{
		// ping
		v1.GET("/site/ping", controllers.Ping)
		// 刷新所有仓库
		v1.POST("/site/refresh/all", controllers.RefreshAll)
		// 刷新特定仓库
		//v1.POST("/site/refresh/repo/:id", controllers.RefreshByRepo)
		// 刷新特定标签
		//v1.POST("/site/refresh/tag/:id", controllers.RefreshByTag)

		// 登录
		//v1.POST("/user/login", controllers.Login)
		// 注册
		//v1.POST("/user/register", controllers.Register)

		// 获取所有仓库
		v1.GET("/repos", controllers.GetRepos)
		// 获取特定仓库下所有标签
		v1.GET("/repos/:id/tags", controllers.GetRootTagInfos)

		// 获取特定标签下文章基本信息列表
		v1.GET("tags/:id/articles", controllers.GetTagView)

		// 获取文章详情
		v1.GET("articles/:id", controllers.GetArticle)
	}

	return r
}
