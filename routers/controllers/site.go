package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/service"
)

func Ping(c *gin.Context) {
	c.JSON(200, serializer.SuccessResponse("", "pong"))
}

// 刷新某 Tag 下的 articles
//func RefreshByTag(c *gin.Context) {
//	c.JSON(200, serializer.Response{
//		Code: 0,
//		Data: "pong",
//	})
//}

// 刷新某 Repo 下的 articles
//func RefreshByRepo(c *gin.Context) {
//	c.JSON(200, serializer.Response{
//		Code: 0,
//		Data: "pong",
//	})
//}

// 刷新所有内容
func RefreshAll(c *gin.Context) {
	res := service.RefreshGlobal()
	c.JSON(200, res)
}
