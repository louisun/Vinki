package services

import (
	"github.com/vinki/db"
)

func GetAllTags() ([]string, error) {
	tags, err := db.GetTags()
	if err != nil {
		log.Errorf("Get All Tags failed: %v", err)
		return nil, err
	}
	return tags, nil
}

func AddTags(tags []*db.Tag) error {
	err := db.AddTags(tags)
	if err != nil {
		log.Errorf("Add tags failed: %v", err)
		return err
	}
	return nil
}
