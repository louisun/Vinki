package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/pkg/utils"
	"github.com/louisun/vinki/service"
)

func Ping(c *gin.Context) {
	c.JSON(200, serializer.CreateSuccessResponse("pong", "服务器状态正常"))
}

// GetSiteConfig 获取用户站点配置
func GetSiteConfig(c *gin.Context) {
	user := GetCurrentUserFromCtx(c)

	if user != nil {
		c.JSON(200, serializer.CreateUserResponse(user))
		return
	}
	c.JSON(200, serializer.GetUnauthorizedResponse())
	return
}

// RefreshAll 刷新所有内容
func RefreshAll(c *gin.Context) {
	res := service.RefreshGlobal()
	c.JSON(200, res)
}

// RefreshByRepo 刷新某 Repo 下的 articles
func RefreshByRepo(c *gin.Context) {
	var repo service.RepoRequest
	if err := c.ShouldBindJSON(&repo); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		utils.Log().Info("repo.RepoName = ", repo.RepoName)
		res := service.RefreshRepo(repo.RepoName)
		c.JSON(200, res)
	}
}

// RefreshByTag 刷新某 Tag 下的 articles
func RefreshByTag(c *gin.Context) {
	var tag service.RepoTagRequest
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := service.RefreshTag(tag.RepoName, tag.TagName)
		c.JSON(200, res)
	}
}

// Search 搜索
func Search(c *gin.Context) {
	res := service.Search(c)
	c.JSON(200, res)
}
