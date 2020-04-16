package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Article struct {
	ID    uint64 `gorm:"primary_key"`
	Title string `gorm:"type:varchar(100);not null"` // 文章标题
	Path  string `gorm:"type:varchar(200);not null"` // 文件路径
	HTML  string `gorm:"type:text"`                  // Markdown 渲染后的 HTML
	TagID uint64 // 标签 ID
	Tag   Tag    `gorm:"PRELOAD:false;save_associations:false"` // 关联的 Tag
}

type ArticleInfo struct {
	ID    uint64
	Title string
}

func (Article) TableName() string {
	return "article"
}

// GetArticleByID 通过 ID 获取 Article
func GetArticleByID(ID uint64) (Article, error) {
	var article Article
	result := DB.First(&article, ID)
	return article, result.Error
}

// GetArticleInfosByTagID 根据 tagID 返回 Article 基本信息
func GetArticleInfosByTagID(tagID uint64) ([]ArticleInfo, error) {
	var articles []ArticleInfo
	result := DB.Model(&Article{}).Where("tag_id = ?", tagID).Select("id, title").Scan(&articles)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return articles, result.Error
	}
	return articles, nil
}

// TruncateArticles 清空 Article 表
func TruncateArticles() error {
	err := DB.Model(&Article{}).Delete(&Article{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteArticlesByRepoID 删除该 Repo 下的 Article
func DeleteArticlesByRepoID(repoID uint64) error {
	result := DB.Preload("Tag", "repoID = ?", repoID).Delete(&Article{})
	return result.Error
}

// DeleteArticlesByTagID 删除该 Tag 下的 Article
func DeleteArticlesByTagID(tagID uint64) error {
	result := DB.Where("tag_id = ?", tagID).Delete(&Article{})
	return result.Error
}

// AddArticle 添加 Article
func AddArticle(article *Article) error {
	err := DB.Create(article).Error
	return err
}

// AddArticles 批量添加 Articles
func AddArticles(articles []*Article) error {
	DB.LogMode(false)
	defer func() {
		if gin.Mode() == gin.TestMode {
			DB.LogMode(true)
		}
	}()
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	for _, article := range articles {
		if err := tx.Create(article).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateArticle 更新文章
func UpdateArticle(id uint64, title string, html string) error {
	result := DB.Model(&Article{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": title,
		"html":  html,
	})
	return result.Error
}
