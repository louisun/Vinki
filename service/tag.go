package service

import (
	"github.com/vinki/db"
)

// 获取文章
func GetTagHtml(tagName string) (*db.Tag, error) {
	tag, err := db.GetTag(tagName)
	if err != nil {
		log.Errorf("Get tag by name failed: %v, tag: %s", err, tag)
		return nil, err
	}
	return tag, nil
}

// 获取所有标签名
func GetAllTags() ([]string, error) {
	tags, err := db.GetTags()
	if err != nil {
		log.Errorf("Get All Tags failed: %v", err)
		return nil, err
	}
	return tags, nil
}

// 添加标签名
func AddTags(tags []*db.Tag) error {
	err := db.AddTags(tags)
	if err != nil {
		log.Errorf("Add tags failed: %v", err)
		return err
	}
	return nil
}
