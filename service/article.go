package service

import (
	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/model"
	"github.com/louisun/vinki/pkg/serializer"
)

// ArticleView 文章视图
type ArticleView struct {
	Title string
	HTML  string
}

// GetArticleDetail 获取文章详情
func GetArticleDetail(repoName string, tagName string, articleName string) serializer.Response {
	article, err := model.GetArticle(repoName, tagName, articleName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.CreateGeneralParamErrorResponse("文章 ID 不存在", err)
		}

		return serializer.CreateDBErrorResponse("", err)
	}

	view := ArticleView{
		Title: article.Title,
		HTML:  article.HTML,
	}

	return serializer.CreateSuccessResponse(view, "")
}
