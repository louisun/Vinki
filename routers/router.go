package routers

import (
	"time"

	"github.com/louisun/vinki/pkg/conf"

	"github.com/louisun/vinki/middleware"

	"github.com/gin-contrib/static"
	"github.com/louisun/vinki/bootstrap"

	"github.com/gin-contrib/cors"
	"github.com/louisun/vinki/controllers"

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

	// API 接口
	v1 := r.Group("/api/v1")

	// Session 中间件，将处理 Cookie 和会话，在 Context 中设置 session
	v1.Use(middleware.Session(conf.GlobalConfig.System.SessionSecret))

	// TODO 跨域相关中间件

	// 中间件：在 Context 中设置当前 Session 对应的用户对象
	v1.Use(middleware.InitCurrentUserIfExists())

	/*
		API 路由
	*/
	{
		// 站点相关路由
		site := v1.Group("site")
		{
			site.GET("ping", controllers.Ping)
			// TODO 验证码
			// TODO 站点全局配置
		}

		// 用户相关路由
		user := v1.Group("user")
		{
			user.POST("login", controllers.UserLogin)
			user.POST("logout", controllers.UserLogout)
			user.POST("", controllers.UserRegister)
			user.PUT("", controllers.UserResetPassword)
			// 封禁用户
			user.POST("ban", middleware.RequireAuth(), middleware.RequireAdmin(), controllers.BanUser)
			// TODO 用户激活
			// TODO OAuth2 登录

		}

		// 要求认证
		auth := v1.Group("")
		// 中间件：要求已登录，即在 Context 中设置了对应的用户对象
		auth.Use(middleware.RequireAuth())

		{
			admin := auth.Group("admin", middleware.RequireAdmin())
			{
				// 刷新所有仓库
				admin.POST("refresh/all", controllers.RefreshAll)
				// 刷新特定仓库
				//auth.POST("/site/refresh/repo/:id", controllers.RefreshByRepo)
				// 刷新特定标签
				//auth.POST("/site/refresh/tag/:id", controllers.RefreshByTag)
				admin.GET("/applications", controllers.GetApplications)
				admin.POST("/application/activate", controllers.ActivateUser)
				admin.POST("/application/reject", controllers.RejectUserApplication)
			}
			// 向管理员申请激活
			auth.POST("/apply", controllers.ApplyForActivate)
			// 搜索内容
			auth.GET("/search", controllers.Search)

			active := auth.Group("", middleware.CheckPermission())
			{
				// 获取所有仓库
				active.GET("/repos", controllers.GetRepos)
				// 获取特定仓库下所有标签
				active.GET("/tags", controllers.GetTopTags)
				// 获取特定标签信息
				active.GET("/tag", controllers.GetTagView)
				// 获取文章详情
				active.GET("/article", controllers.GetArticle)
			}
		}
	}

	return r
}
