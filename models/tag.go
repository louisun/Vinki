package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

// Tag 标签
type Tag struct {
	ID         uint64         `gorm:"primary_key"`
	Path       string         `gorm:"type:varchar(200);unique_index:idx_repo_tag;index:path;not null"` // Tag 目录路径
	RepoID     uint64         `gorm:"unique_index:idx_repo_tag"`                                       // 所属仓库
	Name       string         `gorm:"type:varchar(50);not null"`                                       // 标签名称
	ParentPath sql.NullString `gorm:"type:varchar(200);index:parent_path;default: null"`               // 父标签路径
	Repo       Repo           `gorm:"save_associations:false"`                                         // 关联的 Repo
}

func (Tag) TableName() string {
	return "tag"
}

type TagInfo struct {
	ID   uint   // 标签 ID
	Name string // 标签名
}

// TagView 标签视图
type TagView struct {
	ID      uint64
	Name    string
	SubTags []TagInfo // 子标签列表
}

// GetTagsByRepo 根据 repo 名获取所有 Tag 信息
func GetTagsByRepo(repoID uint64) ([]Tag, error) {
	var tags []Tag
	result := DB.Where("repo_id = ?", repoID).Order("name").Find(&tags)
	return tags, result.Error
}

// GetRootTagsByRepo 根据 repo 名获取所有一级标签的 Tag 信息
func GetRootTagsByRepo(repoID uint64) ([]Tag, error) {
	var tags []Tag
	result := DB.Where("repo_id = ? and parent_path is null", repoID).Order("name").Find(&tags)
	return tags, result.Error
}

// GetRootTagInfosByRepo 根据 repo 名获取所有一级标签的 TagInfo 信息
func GetRootTagInfosByRepo(repoID uint64) ([]TagInfo, error) {
	var tags []TagInfo
	result := DB.Model(&Tag{}).Select("id, name").
		Where("repo_id = ? and parent_path is null", repoID).Order("name").Scan(&tags)
	return tags, result.Error
}

// GetTagViewByID 通过 TagID 获取 TagView
func GetTagViewByID(tagID uint64) (TagView, error) {
	var tagView TagView
	var tag Tag
	// 本标签信息
	result := DB.Where("id = ?", tagID).Select("id, name, path").Find(&tag)
	if result.Error != nil {
		return tagView, result.Error
	}
	// 一级子标签列表
	var subTags []TagInfo
	result = DB.Model(&Tag{}).Where("parent_path = ?", tag.Path).Select("id, name").Scan(&subTags)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tagView, result.Error
	}
	tagView.ID = tag.ID
	tagView.Name = tag.Name
	tagView.SubTags = subTags
	return tagView, nil
}

// GetFlatTagViewByID 通过 TagID 获取平铺的 TagView
func GetFlatTagViewByID(tagID uint64) (TagView, error) {
	var list []TagInfo
	tagView, err := GetTagViewByID(tagID)
	if err != nil {
		return tagView, err
	}
	for _, tagInfo := range tagView.SubTags {
		traverseTagInfo(&tagInfo, list)
	}
	// 替换 SubTags 为所有子标签
	tagView.SubTags = list
	return tagView, nil
}

func traverseTagInfo(tagInfo *TagInfo, list []TagInfo) {
	var subTags []TagInfo
	result := DB.Model(&Tag{}).Where("parent_path = ?", tagInfo.ID).Select("id, name").Scan(&subTags)
	// 无子标签也立即返回
	if result.Error != nil {
		return
	}
	list = append(list, subTags...)
	for _, subTagInfo := range subTags {
		traverseTagInfo(&subTagInfo, list)
	}
}

// TruncateTags 清空 Tag 表
func TruncateTags() error {
	result := DB.Delete(&Tag{})
	return result.Error
}

// TruncateTagsByRepo 清空某个 repo 下的 Tags
func DeleteTagsByRepo(repoID uint64) error {
	result := DB.Where("repo_id = ?", repoID).Delete(&Tag{})
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
