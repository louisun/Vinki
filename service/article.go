package service

import (
	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
)

type ArticleView struct {
	Title string
	HTML  string
}

// 获取文章详情
func GetArticleDetail(repoName string, tagName string, articleName string) serializer.Response {
	article, err := models.GetArticle(repoName, tagName, articleName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("文章 ID 不存在", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	view := ArticleView{
		Title: article.Title,
		HTML:  article.HTML,
	}
	return serializer.SuccessResponse(view, "")
}

// 获取某 Tag 下的文章列表
func GetArticleList(repoName string, tagName string) serializer.Response {
	articles, err := models.GetArticleList(repoName, tagName)
	if err != nil {
		return serializer.DBErrorResponse("", err)
	}
	return serializer.SuccessResponse(articles, "")
}

// 批量添加 Articles
func addArticles(articles []*models.Article) error {
	err := models.AddArticles(articles)
	return err
}

// deleteArticlesByTag 删除某 Tag 下的 Article
func deleteArticlesByTag(repoName string, tagName string) error {
	err := models.DeleteArticlesByTagName(repoName, tagName)
	return err
}

// deleteArticlesByRepo 删除某 Repo 下的 Article
func deleteArticlesByRepoName(repoName string) error {
	err := models.DeleteArticlesByRepoName(repoName)
	return err
}

// TruncateArticles 清空所有 Articles
func truncateArticles() error {
	err := models.TruncateArticles()
	return err
}
