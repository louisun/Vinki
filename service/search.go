package service

import (
	"errors"

	"github.com/louisun/vinki/models"

	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/pkg/serializer"
)

const (
	typeTag         = "tag"
	typeArticleName = "article"
)

// Search 搜索
func Search(c *gin.Context) serializer.Response {
	searchType := c.Query("type")
	repoName := c.Query("repo")
	keyword := c.Query("keyword")
	if keyword == "" {
		return serializer.CreateGeneralParamErrorResponse("", errors.New("搜索关键词不能为空"))
	}
	switch searchType {
	case typeTag:
		tags, err := models.GetTagsBySearchName(repoName, keyword)
		if err != nil {
			return serializer.CreateDBErrorResponse("", err)
		}
		return serializer.CreateSuccessResponse(tags, "")
	case typeArticleName:
		articleTagInfos, err := models.GetArticlesBySearchParam(repoName, keyword)
		if err != nil {
			return serializer.CreateDBErrorResponse("", err)
		}
		return serializer.CreateSuccessResponse(articleTagInfos, "")
	default:
		return serializer.CreateGeneralParamErrorResponse("", errors.New("搜索类型不正确"))
	}
}
