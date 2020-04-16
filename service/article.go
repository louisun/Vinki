package service

import (
	"github.com/jinzhu/gorm"
	"github.com/vinki/models"
	"github.com/vinki/pkg/serializer"
)

type ArticleView struct {
	ID    uint64
	Title string
	HTML  string
}

// 获取文章详情
func GetArticleDetail(articleID uint64) serializer.Response {
	article, err := models.GetArticleByID(articleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("文章 ID 不存在", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	view := ArticleView{
		ID:    article.ID,
		Title: article.Title,
		HTML:  article.HTML,
	}
	return serializer.SuccessResponse(view, "")
}

// 获取某 Tag 下的文章列表
func GetArticleListByTagID(tagID uint64) serializer.Response {
	articles, err := models.GetArticleInfosByTagID(tagID)
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
func deleteArticlesByTag(tagID uint64) error {
	err := models.DeleteArticlesByTagID(tagID)
	return err
}

// deleteArticlesByRepo 删除某 Repo 下的 Article
func deleteArticlesByRepo(repoID uint64) error {
	err := models.DeleteArticlesByRepoID(repoID)
	return err
}

// TruncateArticles 清空所有 Articles
func truncateArticles() error {
	err := models.TruncateArticles()
	return err
}
