package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vinki/pkg/serializer"
	"github.com/vinki/service"
)

// 获取某文章详细信息
func GetArticle(c *gin.Context) {
	var (
		id  uint64
		err error
	)
	s := c.Param("id")
	if id, err = strconv.ParseUint(s, 10, 64); err != nil {
		c.JSON(200, serializer.ParamErrorResponse("article id 错误", err))
	}
	res := service.GetArticleDetail(id)
	c.JSON(200, res)
}

// 获取某 Tag 的文章基本信息
func GetTagView(c *gin.Context) {
	var (
		id  uint64
		err error
	)
	s := c.Param("id")
	if id, err = strconv.ParseUint(s, 10, 64); err != nil {
		c.JSON(200, serializer.ParamErrorResponse("tag id 错误", err))
	}
	res := service.GetTagViewByID(id)
	c.JSON(200, res)
}

// 获取某 Repo 下的标签列表
func GetRootTagInfos(c *gin.Context) {
	var (
		id  uint64
		err error
	)
	s := c.Param("id")
	if id, err = strconv.ParseUint(s, 10, 64); err != nil {
		c.JSON(200, serializer.ParamErrorResponse("repo id 错误", err))
	}
	res := service.GetRootTagInfosByRepo(id)
	c.JSON(200, res)
}

// 获取所有 Repo 列表
func GetRepos(c *gin.Context) {
	res := service.GetRepoInfos()
	c.JSON(200, res)
}
