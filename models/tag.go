package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

// Tag 标签
type Tag struct {
	ID         uint64         `gorm:"primary_key"`
	Path       string         `gorm:"type:varchar(200);index:path;not null"`             // Tag 目录路径
	Name       string         `gorm:"type:varchar(200);index:name;not null"`             // 标签名称（多级拼合）
	ParentPath sql.NullString `gorm:"type:varchar(200);index:parent_path;default: null"` // 父标签路径
	RepoName   string
	Repo       Repo `gorm:"foreignkey:RepoName;association_foreignkey:Name;PRELOAD:false;save_associations:false"`
}

func (Tag) TableName() string {
	return "tag"
}

// TagView 标签视图
type TagView struct {
	Name    string
	SubTags []string // 子标签列表
}

// GetTagsByRepo 根据 repo 名获取所有 Tag 信息
func GetTagsByRepoName(repoName string) ([]Tag, error) {
	var tags []Tag
	result := DB.Where("repo_name = ?", repoName).Order("repoName").Find(&tags)
	return tags, result.Error
}

func GetTag(repoName, tagName string) (Tag, error) {
	var tag Tag
	result := DB.Where("repo_name = ? and name = ?", repoName, tagName).First(&tag)
	return tag, result.Error
}

// GetTagsBySearchName 根据 repo 名和 tagName 搜索 Tags
func GetTagsBySearchName(repoName, tagName string) ([]string, error) {
	var tags []string
	pattern := "%" + tagName + "%"
	result := DB.Model(&Tag{}).Where("repo_name = ? AND name LIKE ?", repoName, pattern).Order("length(`name`)").Pluck("name", &tags)
	return tags, result.Error
}

// GetRootTagsByRepo 根据 repo 名获取所有一级标签的 Tag 信息
func GetRootTagsByRepo(repoName string) ([]Tag, error) {
	var tags []Tag
	result := DB.Where("repo_name = ? AND parent_path IS NULL", repoName).Order("name").Find(&tags)
	return tags, result.Error
}

// GetTopTagInfosByRepo 根据 repo 名获取所有一级标签的名称
func GetTopTagInfosByRepo(repoName string) ([]string, error) {
	tags := make([]string, 0)
	result := DB.Model(&Tag{}).Where("repo_name = ? AND parent_path IS NULL", repoName).
		Order("name").Pluck("name", &tags)
	return tags, result.Error
}

// GetTagView 通过 repoName 和 tagName 获取 TagView
func GetTagView(repoName string, tagName string) (TagView, error) {
	var tagView TagView
	var tag Tag
	// 本标签信息
	result := DB.Where("repo_name = ? AND name = ?", repoName, tagName).Select("name, path").Find(&tag)
	if result.Error != nil {
		return tagView, result.Error
	}
	// 一级子标签列表
	subTags := make([]string, 0)
	result = DB.Model(&Tag{}).Where("parent_path = ?", tag.Path).Pluck("name", &subTags)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tagView, result.Error
	}
	tagView.Name = tag.Name
	tagView.SubTags = subTags
	return tagView, nil
}

// GetFlatTagView 通过 TagID 获取平铺的 TagView
func GetFlatTagView(repoName string, tagName string) (TagView, error) {
	var tag Tag
	var list []string
	var tagView TagView
	// 本标签信息
	result := DB.Where("repo_name = ? AND name = ?", repoName, tagName).Select("name, path").Find(&tag)
	if result.Error != nil {
		return tagView, result.Error
	}
	// 一级子标签列表
	var subTags []Tag
	result = DB.Model(&Tag{}).Where("parent_path = ?", tag.Path).Select("name, path").Find(&subTags)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tagView, result.Error
	}
	// 多级子标签列表
	for _, subTag := range subTags {
		list = append(list, subTag.Name)
		traverseTagInfo(&list, &subTag)
	}
	tagView.Name = tag.Name
	tagView.SubTags = list
	return tagView, nil
}

func traverseTagInfo(list *[]string, parentTag *Tag) {
	var subTags []Tag
	result := DB.Model(&Tag{}).Where("parent_path = ?", parentTag.Path).Select("name, path").Find(&subTags)
	// 无子标签也立即返回
	if result.Error != nil {
		return
	}
	for _, tag := range subTags {
		*list = append(*list, tag.Name)
	}
	for _, tag := range subTags {
		traverseTagInfo(list, &tag)
	}
}

// TruncateTags 清空 Tag 表
func TruncateTags() error {
	result := DB.Delete(&Tag{})
	return result.Error
}

// DeleteTagsByRepo 清空某个 repo 下的 Tags
func DeleteTagsByRepo(repoName string) error {
	result := DB.Where("repo_name = ?", repoName).Delete(&Tag{})
	return result.Error
}

// DeleteTag 清空某个 repo 下的 tag
func DeleteTag(repoName string, tagName string) error {
	result := DB.Where("repo_name = ? and tag = ?", repoName, tagName).Delete(&Tag{})
	return result.Error
}

// AddTag 添加 Tag
func AddTag(tag Tag) error {
	err := DB.Create(&tag).Error
	return err
}

// AddTags 批量添加 Tags
func AddTags(tags []*Tag) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	for _, tag := range tags {
		if err := tx.Create(tag).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
