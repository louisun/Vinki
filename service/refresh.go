package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/louisun/vinki/pkg/serializer"

	"github.com/louisun/vinki/pkg/conf"

	"github.com/louisun/vinki/models"

	"github.com/louisun/vinki/pkg/utils"
)

// RefreshGlobal 全局刷新
func RefreshGlobal() serializer.Response {
	// 1. 清空数据库
	if err := truncateDB(); err != nil {
		return serializer.DBErrorResponse("", err)
	}
	// 遍历每个 Repo 配置项
	for _, r := range conf.GlobalConfig.Repositories {
		if utils.ExistsDir(r.Root) {
			var repoPath = r.Root
			var tagPath2Name = make(map[string]string)
			if strings.HasSuffix(repoPath, "/") {
				repoPath = strings.TrimSuffix(repoPath, "/")
			}
			repo := models.Repo{
				Name: filepath.Base(repoPath),
				Path: repoPath,
			}
			// 2. 创建 Repo
			if err := addRepo(&repo); err != nil {
				utils.Log().Errorf("addRepo failed: %v", err)
				return serializer.DBErrorResponse("", err)
			}
			// 遍历 Repo，找到 Tag 和 Article 的路径
			tagPaths, t2f, err := traverseRepo(r)
			tags := make([]*models.Tag, 0, len(tagPaths))
			for _, tp := range tagPaths {
				parentPath := filepath.Dir(tp)
				tagName := strings.Join(strings.Split(strings.TrimPrefix(tp, repo.Path+"/"), "/"), "--")
				if parentPath == repoPath {
					// 一级目录
					tags = append(tags, &models.Tag{
						Path:     tp,
						RepoName: repo.Name,
						Name:     tagName,
					})
				} else {
					// 子目录
					tags = append(tags, &models.Tag{
						Path:     tp,
						RepoName: repo.Name,
						Name:     tagName,
						ParentPath: sql.NullString{
							String: parentPath,
							Valid:  true,
						},
					})
				}
			}
			// 3. 创建 Tags
			if err = AddTags(tags); err != nil {
				utils.Log().Errorf("AddTags failed: %v", err)
				return serializer.DBErrorResponse("", err)
			}
			for _, tag := range tags {
				tagPath2Name[tag.Path] = tag.Name
			}
			// 4. 创建 Articles
			for tagPath, fileInfos := range t2f {
				articleList := make([]*models.Article, 0, len(fileInfos))
				for _, fileInfo := range fileInfos {
					var htmlBytes []byte
					htmlBytes, err = utils.RenderMarkdown(fileInfo.Path)
					if err != nil {
						utils.Log().Errorf("Generate wiki html failed, err: %v", err)
						return serializer.InternalErrorResponse("", err)
					}
					articleList = append(articleList, &models.Article{
						RepoName: repo.Name,
						TagName:  tagPath2Name[tagPath],
						Title:    fileInfo.BriefName,
						Path:     fileInfo.Path,
						HTML:     string(htmlBytes),
					})
				}
				err = addArticles(articleList)
				if err != nil {
					utils.Log().Errorf("addArticles failed: %v", err)
					return serializer.DBErrorResponse("", err)
				}
			}
		} else {
			err := errors.New(fmt.Sprintf("Repo path does not exist: %s", r.Root))
			utils.Log().Error(err)
			return serializer.InternalErrorResponse("", err)
		}
	}
	utils.Log().Info("Refresh database success!")
	return serializer.SuccessResponse("", "同步仓库成功")
}

// truncateDB 清空所有 Article、Tag、Repo
func truncateDB() error {
	// 清空 article 数据库
	err := truncateArticles()
	if err != nil {
		utils.Log().Errorf("Truncate article db failed: %v", err)
	}
	// 清空 tag 数据库
	err = truncateTags()
	if err != nil {
		utils.Log().Errorf("Truncate tag db failed: %v", err)
		return err
	}
	// 清空 repo 数据库
	err = truncateRepo()
	if err != nil {
		utils.Log().Errorf("Truncate repo db failed: %v", err)
		return err
	}
	return nil
}

// traverseRepo 遍历 Repo 目录，生成标签与文档的映射树
func traverseRepo(repo conf.DirectoryConfig) (tagPaths []string, t2f utils.TagPath2FileInfo, err error) {
	t2f = make(map[string][]*utils.FileInfo)
	var root = repo.Root
	var exclude = repo.Exclude

	// rootPrefix 用于拆分标签
	var rootPrefix = root
	if !strings.HasSuffix(rootPrefix, "/") {
		rootPrefix = rootPrefix + "/"
	}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// 目录 -> Tag
			// 跳过 Exclude 配置项指定的目录
			if strings.HasPrefix(info.Name(), ".") || utils.IsDirectoryInList(info.Name(), exclude) {
				return filepath.SkipDir
			}
			// 保存 Tag path
			tagPaths = append(tagPaths, path)
		} else if strings.HasSuffix(info.Name(), ".md") {
			// md 文件 -> Article
			// 保存 Tag 路径 -> FileInfo
			dir := filepath.Dir(path)
			if paths, ok := t2f[dir]; ok {
				paths = append(paths, &utils.FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				t2f[dir] = paths
			} else {
				paths = make([]*utils.FileInfo, 0, 10)
				paths = append(paths, &utils.FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				t2f[dir] = paths
			}
		}
		return nil
	})
	return
}
