package routers

import (
	"time"

	"github.com/gin-contrib/static"
	"github.com/louisun/vinki/bootstrap"

	"github.com/gin-contrib/cors"
	"github.com/louisun/vinki/routers/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	if gin.Mode() == gin.TestMode {
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	r.Use(static.Serve("/", bootstrap.StaticFS))

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
