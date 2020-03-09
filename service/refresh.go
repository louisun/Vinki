package service

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/vinki/utils"

	"github.com/vinki/db"
)

var numberRegex = regexp.MustCompile(`^(\d+)\..*`)

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
	}
	sort.Strings(tagNames)
	var articleCount = 1
	// 4. 生成主页
	htmlContent, err := utils.GenerateWikiHtml(utils.Config.Custom.Home, 1, "Home", tagNames, []string{})
	if err != nil {
		log.Errorf("Generate home html failed, err: %v", err)
		return err
	}
	err = AddArticles([]*db.Article{{
		ID:          articleCount,
		Title:       "Vinki",
		FilePath:    utils.Config.Custom.Home,
		Tag:         "Home",
		HtmlContent: htmlContent,
	}})
	articleCount++
	if err != nil {
		return err
	}

	// 5. 生成文章
	for tagName, fileInfos := range utils.Tag2FilePath {
		articleList := make([]*db.Article, 0, len(fileInfos))
		files := make([]string, 0, len(fileInfos))
		orderNumbers := make([]int, 0, len(fileInfos))
		order2Name := make(map[int]string, len(fileInfos))
		for _, fileInfo := range fileInfos {
			match := numberRegex.FindStringSubmatch(fileInfo.BriefName)
			if match != nil {
				n, err := strconv.Atoi(match[1])
				if err != nil {
					// 转换出错，添加到列表头部
					files = append([]string{fileInfo.BriefName}, files...)
					continue
				}
				order2Name[n] = fileInfo.BriefName
				orderNumbers = append(orderNumbers, n)
			} else {
				// 非数字编号先添加到列表
				files = append(files, fileInfo.BriefName)
			}
		}
		sort.Strings(files)
		sort.Sort(sort.Reverse(sort.IntSlice(orderNumbers)))
		for _, n := range orderNumbers {
			// 数字编号，添加到列表头部
			files = append([]string{order2Name[n]}, files...)
		}

		// 5.1 为每个 tag 生成完整的 Html 页面
		htmlContent, err := utils.GenerateWikiHtml(utils.Config.Custom.Tag, len(fileInfos), tagName, tagNames, files)
		if err != nil {
			log.Errorf("Generate tag html failed, err: %v", err)
			return err
		}
		tags = append(tags, &db.Tag{
			ID:          tagCount,
			Name:        tagName,
			HtmlContent: htmlContent,
		})
		tagCount++
		// 5.2 为每篇文章生成完整的 Html
		for _, fileInfo := range fileInfos {
			htmlContent, err := utils.GenerateWikiHtml(fileInfo.Path, len(fileInfos), tagName, tagNames, files)
			if err != nil {
				log.Errorf("Generate wiki html failed, err: %v", err)
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
	// 6. 添加标签
	err = AddTags(tags)
	if err != nil {
		return err
	}
	log.Info("Refresh database success!")
	return nil
}
