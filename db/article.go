package db

import "github.com/jinzhu/gorm"

// 文章
type Article struct {
	ID          int    `gorm:"column:id;primary_key" json:"id"`
	Title       string `gorm:"column:title" json:"title"`               // 文章标题
	FilePath    string `gorm:"column:file_path" json:"file_path"`       // 文件路径
	Tag         string `gorm:"column:tag" json:"tag"`                   // 标签名
	HtmlContent string `gorm:"column:html_content" json:"html_content"` // 完整的渲染后 html 内容
}

// 将 User 的表名设置为 `profiles`
func (Article) TableName() string {
	return "Article"
}

// 根据 Tag 和 Title 获取文章信息
func GetArticleByTagAndTitle(tag string, title string) (*Article, error) {
	var (
		article Article
		err     error
	)
	err = db.Model(&Article{}).Where("tag = ? AND title = ?", tag, title).Find(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// 获取同一 Tag 下的文章名列表
func GetArticleListByTag(tag string) ([]string, error) {
	var articleNames []string
	err := db.Model(&Article{}).Where("tag = ?", tag).Order("title").Pluck("title", &articleNames).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articleNames, nil
}

// 清空 Article 表
func TruncateArticle() error {
	err := db.Model(&Article{}).Delete(&Article{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 添加 Article
func AddArticle(article *Article) error {
	err := db.Model(&Article{}).Create(article).Error
	return err
}
