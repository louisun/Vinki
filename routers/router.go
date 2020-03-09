package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinki/routers/api"
)

func LoadRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	wg := r.Group("/")
	{
		// 静态资源
		r.StaticFS("/static", http.Dir("./static"))
		// 检测服务状态
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 主页
		r.GET("/", api.GetHomePage)
		// 刷新数据库
		r.GET("/refresh", api.Refresh)
	}

	// 缓存-中间件
	wg = r.Group("/wiki")
	{
		// 主页
		wg.GET("/", api.GetHomePage)

		// Wiki 页
		wg.GET("/tags/:tag/:title", api.GetWikiPage)

		// Tag 主页
		wg.GET("/tags/:tag", api.GetTagPage)
	}

	return r
}
