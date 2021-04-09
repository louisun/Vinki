package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/service"
)

// 获取某文章详细信息
func GetArticle(c *gin.Context) {
	repoName := c.Query("repoName")
	tagName := c.Query("tagName")

	articleName := c.Query("articleName")
	if repoName == "" || tagName == "" || articleName == "" {
		c.JSON(200, serializer.CreateGeneralParamErrorResponse("", errors.New("仓库名、标签名和文章名不能为空")))
		return
	}

	res := service.GetArticleDetail(repoName, tagName, articleName)
	c.JSON(200, res)

	return
}

// GetTagView 获取某 Tag 的基本信息
func GetTagView(c *gin.Context) {
	repoName := c.Query("repoName")
	tagName := c.Query("tagName")

	if repoName == "" || tagName == "" {
		c.JSON(200, serializer.CreateGeneralParamErrorResponse("", errors.New("仓库名和标签名不能为空")))
		return
	}

	var flat bool
	if c.Query("flat") == "true" {
		flat = true
	}

	var res serializer.Response

	res = service.GetTagArticleView(repoName, tagName, flat)
	c.JSON(200, res)

	return
}

// GetTopTags 获取某 Repo 下的一级标签列表
func GetTopTags(c *gin.Context) {
	repoName := c.Query("repoName")
	if repoName == "" {
		c.JSON(200, serializer.CreateGeneralParamErrorResponse("", errors.New("仓库名不能为空")))
		return
	}

	res := service.GetTopTagInfosByRepo(repoName)
	c.JSON(200, res)

	return
}

// GetRepos 获取所有 Repo 列表
func GetRepos(c *gin.Context) {
	user := GetCurrentUserFromCtx(c)
	if user != nil {
		c.JSON(200, serializer.CreateSuccessResponse(user.RepoNames, ""))
		return
	}

	c.JSON(200, serializer.GetUnauthorizedResponse())

	return
}
