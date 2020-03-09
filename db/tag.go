package db

import "github.com/jinzhu/gorm"

// 默认标签名就是目录名，但可修改
type Tag struct {
	ID          int    `gorm:"column:id;primary_key" json:"id"`
	Name        string `gorm:"column:name;unique;not null" json:"name"`         // 标签名
	HtmlContent string `gorm:"column:html_content;unique;not null" json:"name"` // Html 页面
}

// 将 User 的表名设置为 `profiles`
func (Tag) TableName() string {
	return "Tag"
}

// 获取标签信息
func GetTags() ([]string, error) {
	var (
		tags []string
		err  error
	)
	err = db.Model(Tag{}).Order("name").Pluck("name", &tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

// 清空 Tag 表
func TruncateTags() error {
	err := db.Model(&Tag{}).Delete(&Tag{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 添加 Tag
func AddTags(tags []*Tag) error {
	var err error
	for _, tag := range tags {
		err = db.Model(Tag{}).Create(tag).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// 更新 Tag 的 Html
func UpdateTagHtml(tag *Tag) error {
	err := db.Model(&Tag{}).Where("name = ?", tag.Name).Update("html_content", tag.HtmlContent).Error
	return err
}

// 获取 Tag Html
func GetTag(tagName string) (*Tag, error) {
	var (
		tag Tag
		err error
	)
	err = db.Model(&Article{}).Where("name = ?", tagName).Find(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
