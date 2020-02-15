package services

import (
	"sort"

	"github.com/vinki/pkg/utils"

	"github.com/vinki/db"
)

var log = utils.DefaultLog

// 清空数据库
func TruncateDB() error {
	// 清空Article数据库
	err := db.TruncateArticle()
	if err != nil {
		log.Errorf("Truncate article db failed: %v", err)
		return err
	}
	// 清空tag数据库
	err = db.TruncateTags()
	if err != nil {
		log.Errorf("Truncate tag db failed: %v", err)
		return err
	}
	return nil
}

// 刷新数据库
func Refresh() error {
	// 1. 遍历目录，生成树状结构
	err := utils.Traverse()
	if err != nil {
		log.Errorf("Traverse root directories failed: %v", err)
		return err
	}
	// 2. 清空数据库
	err = TruncateDB()
	if err != nil {
		return err
	}
	// 3. 添加标签
	tagNames := make([]string, 0, len(utils.Tag2FilePath))
	tags := make([]*db.Tag, 0, len(utils.Tag2FilePath))
	var tagCount = 1
	for tagName := range utils.Tag2FilePath {
		tagNames = append(tagNames, tagName)
		tags = append(tags, &db.Tag{
			ID:   tagCount,
			Name: tagName,
		})
		tagCount++
	}
	err = AddTags(tags)
	if err != nil {
		return err
	}
	var articleCount = 1
	// 4. 生成文章
	for tagName, fileInfos := range utils.Tag2FilePath {
		articleList := make([]*db.Article, 0, len(fileInfos))
		files := make([]string, 0, len(fileInfos))
		for _, fileInfo := range fileInfos {
			files = append(files, fileInfo.BriefName)
		}
		sort.Strings(files)
		for _, fileInfo := range fileInfos {
			// 为每篇文章生成完整的 Html
			htmlContent, err := utils.GenerateWikiHtml(fileInfo.Path, len(fileInfos), tagName, tagNames, files)
			if err != nil {
				log.Errorf("Generate Wiki Html failed, err: %v", err)
				return err
			}
			articleList = append(articleList, &db.Article{
				ID:          articleCount,
				Title:       fileInfo.BriefName,
				FilePath:    fileInfo.Path,
				Tag:         tagName,
				HtmlContent: htmlContent,
			})
			articleCount++
		}
		// 添加一组文章（同一个标签下）
		err = AddArticles(articleList)
		if err != nil {
			return err
		}
	}
	// 5. 为主页、标签主页生成 HTML
	log.Info("Refresh database success!")
	return nil
}
