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

const maxCorsHour = 12

func handleCors(r *gin.Engine) {
	if gin.Mode() == gin.TestMode {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           maxCorsHour * time.Hour,
		}))
	}
}

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	handleCors(r)
	r.Use(static.Serve("/", bootstrap.StaticFS))

	v1 := r.Group("/api/v1")

	v1.Use(middleware.Session(conf.GlobalConfig.System.SessionSecret))

	v1.Use(middleware.InitCurrentUserIfExists())

	{
		site := v1.Group("site")
		{
			site.GET("ping", controllers.Ping)
			site.GET("config", middleware.RequireAuth(), controllers.GetSiteConfig)
		}

		user := v1.Group("user")
		{
			user.POST("login", controllers.UserLogin)
			user.POST("logout", controllers.UserLogout)
			user.POST("", controllers.UserRegister)
			user.PUT("", controllers.UserResetPassword)
			user.POST("ban", middleware.RequireAuth(), middleware.RequireAdmin(), controllers.BanUser)
		}

		auth := v1.Group("")
		auth.Use(middleware.RequireAuth())

		{
			admin := auth.Group("admin", middleware.RequireAdmin())
			{
				admin.POST("refresh/all", controllers.RefreshAll)
				admin.POST("refresh/repo", controllers.RefreshByRepo)
				admin.POST("refresh/tag", controllers.RefreshByTag)
				admin.GET("applications", controllers.GetApplications)
				admin.POST("application/activate", controllers.ActivateUser)
				admin.POST("application/reject", controllers.RejectUserApplication)
				admin.GET("config/repo", controllers.GetCurrentRepo)
				admin.POST("config/repo", controllers.SetCurrentRepo)
			}
			auth.POST("apply", controllers.ApplyForActivate)
			auth.GET("search", controllers.Search)

			active := auth.Group("", middleware.CheckPermission())
			{
				active.GET("repos", controllers.GetRepos)
				active.GET("tags", controllers.GetTopTags)
				active.GET("tag", controllers.GetTagView)
				active.GET("article", controllers.GetArticle)
			}
		}
	}

	return r
}
