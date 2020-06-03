package models

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var numberRegex = regexp.MustCompile(`^(\d+)\..*`)

type Article struct {
	ID       uint64 `gorm:"primary_key"`
	Title    string `gorm:"type:varchar(100);index:title;not null"` // 文章标题
	Path     string `gorm:"type:varchar(200);not null"`             // 文件路径
	HTML     string `gorm:"type:text"`                              // Markdown 渲染后的 HTML
	TagName  string
	RepoName string
	Tag      Tag  `gorm:"foreignkey:TagName;association_foreignkey:Name;PRELOAD:false;save_associations:false"`  // 关联的 Tag
	Repo     Repo `gorm:"foreignkey:RepoName;association_foreignkey:Name;PRELOAD:false;save_associations:false"` // 关联的 Repo
}

type ArticleTagInfo struct {
	TagName     string `gorm:"column:tag_name" json:"tag"`
	ArticleName string `gorm:"column:title" json:"article"`
}

type Articles []string

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	matchA := numberRegex.FindStringSubmatch(a[i])
	matchB := numberRegex.FindStringSubmatch(a[j])
	var nA = -1
	var nB = -1
	if matchA != nil {
		nA, _ = strconv.Atoi(matchA[1])
	}
	if matchB != nil {
		nB, _ = strconv.Atoi(matchB[1])
	}
	if nA != -1 && nB != -1 {
		return nA < nB
	} else if nA != -1 && nB == -1 {
		return true
	} else if nA == -1 && nB != -1 {
		return false
	} else {
		return a[i] < a[j]
	}
}

func (Article) TableName() string {
	return "article"
}

// GetArticle 通过仓库名、标签名、文章名获取 Article
func GetArticle(repoName string, tagName string, articleName string) (Article, error) {
	var article Article
	result := DB.Where("repo_name = ? AND tag_name = ? AND title = ?", repoName, tagName, articleName).First(&article)
	return article, result.Error
}

// GetArticlesBySearchParam 通过仓库名、文章名搜索 Articles
func GetArticlesBySearchParam(repoName string, articleName string) ([]ArticleTagInfo, error) {
	var articles []ArticleTagInfo
	pattern := "%" + articleName + "%"
	result := DB.Model(&Article{}).Where("repo_name = ? AND title LIKE ?", repoName, pattern).
		Select("title, tag_name").Order("`title`, length(`title`)").Scan(&articles)
	return articles, result.Error
}

// GetArticleList 根据仓库名和标签名获取 Article 列表信息
func GetArticleList(repoName string, tagName string) ([]string, error) {
	articles := make([]string, 0)
	result := DB.Model(&Article{}).Where("repo_name = ? AND tag_name = ?", repoName, tagName).
		Pluck("title", &articles)
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

// DeleteArticlesByRepoName 删除该 Repo 下的 Article
func DeleteArticlesByRepoName(repoName string) error {
	result := DB.Where("repo_name = ?", repoName).Delete(&Article{})
	return result.Error
}

// DeleteArticlesByTagID 删除该 Tag 下的 Article
func DeleteArticlesByTagName(repoName string, tagName string) error {
	result := DB.Where("repo_name = ? AND tag_name = ?", repoName, tagName).Delete(&Article{})
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
