package service

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"

	"github.com/louisun/vinki/pkg/serializer"

	"github.com/louisun/vinki/pkg/conf"

	"github.com/louisun/vinki/models"

	"github.com/louisun/vinki/pkg/utils"
)

type RepoRequest struct {
	RepoName string `json:"repoName" binding:"required"`
}

type RepoTagRequest struct {
	RepoName string `form:"repoName" json:"repoName" binding:"required"`
	TagName  string `form:"tagName" json:"tagName" binding:"required"`
}

type articleTask struct {
	FileInfo *utils.FileInfo
	RepoName string
	TagName  string
}

const (
	ConcurrentSize = 50
	TagSeparator   = "|"
)

// handleArticleTask 处理一个标签下的 Articles
func handleArticleTask(articleTask articleTask, articleChan chan *models.Article) {
	htmlBytes, err := utils.RenderMarkdown(articleTask.FileInfo.Path)
	if err != nil {
		w := fmt.Errorf("RenderMarkdown failed: %w", err)
		utils.Log().Errorf("%v", w)
	}
	article := models.Article{
		Title:    articleTask.FileInfo.BriefName,
		Path:     articleTask.FileInfo.Path,
		HTML:     string(htmlBytes),
		TagName:  articleTask.TagName,
		RepoName: articleTask.RepoName,
	}
	articleChan <- &article
}

// loadLocalRepo 加载本地仓库数据到数据库
func loadLocalRepo(r conf.DirectoryConfig) error {
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
		// 创建 repo
		if err := models.AddRepo(&repo); err != nil {
			w := fmt.Errorf("addRepo failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
		// 遍历 repo，获取标签路径列表、标签对应的文件列表字典
		tagPathList, tagPath2FileListMap, err := traverseRepo(r)

		if err != nil {
			w := fmt.Errorf("traverseRepo failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}

		tags := make([]*models.Tag, 0, len(tagPathList))
		// 根据 tag 路径列表构造 Tag
		for _, tp := range tagPathList {
			parentpath := filepath.Dir(tp)
			tagname := strings.Join(strings.Split(strings.TrimPrefix(tp, repo.Path+"/"), "/"), TagSeparator)

			if parentpath == repoPath {
				// 一级目录
				tags = append(tags, &models.Tag{
					Path:     tp,
					RepoName: repo.Name,
					Name:     tagname,
				})
			} else {
				// 子目录
				tags = append(tags, &models.Tag{
					Path:     tp,
					RepoName: repo.Name,
					Name:     tagname,
					ParentPath: sql.NullString{
						String: parentpath,
						Valid:  true,
					},
				})
			}
		}
		// 3. 创建 Tag
		if err := models.AddTags(tags); err != nil {
			w := fmt.Errorf("addTags failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
		// 保存 tag 路径和标签名的映射关系
		for _, tag := range tags {
			tagPath2Name[tag.Path] = tag.Name
		}

		p, err := ants.NewPool(ConcurrentSize)
		if err != nil {
			w := fmt.Errorf("create goroutine pool failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
		defer p.Release()
		var tasks []articleTask
		// 4. 构造 Article 并创建 articles: 遍历 tagPath2FileListMap
		for tagPath, fileInfos := range tagPath2FileListMap {
			for _, fileInfo := range fileInfos {
				tasks = append(tasks, articleTask{
					fileInfo,
					repo.Name,
					tagPath2Name[tagPath],
				})
			}
		}
		var wg sync.WaitGroup
		var articleChan = make(chan *models.Article, len(tasks))
		var articles = make([]*models.Article, 0, len(tasks))
		for _, task := range tasks {
			wg.Add(1)
			t := task
			_ = p.Submit(func() {
				defer wg.Done()
				handleArticleTask(t, articleChan)
			})
		}
		wg.Wait()
		close(articleChan)
		for article := range articleChan {
			articles = append(articles, article)
		}
		err = models.AddArticles(articles)
		if err != nil {
			w := fmt.Errorf("AddArticles failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
	} else {
		err := fmt.Errorf("repo path not exist: %s", r.Root)
		utils.Log().Error(err)
		return err
	}

	return nil
}

// loadLocalTag 加载本地标签的文章数据到数据库
func loadLocalTag(tag models.Tag) error {
	fileInfos, err := getArticlesInTag(tag)
	if err != nil {
		w := fmt.Errorf("getArticlesInTag failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	articles := make([]*models.Article, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		var htmlBytes []byte
		// Markdown 渲染
		htmlBytes, err := utils.RenderMarkdown(fileInfo.Path)
		if err != nil {
			w := fmt.Errorf("RenderMarkdown failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
		articles = append(articles, &models.Article{
			RepoName: tag.RepoName,
			TagName:  tag.Name,
			Title:    fileInfo.BriefName,
			Path:     fileInfo.Path,
			HTML:     string(htmlBytes),
		})
	}
	err = models.AddArticles(articles)
	if err != nil {
		w := fmt.Errorf("addArticles failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	return nil
}

// RefreshGlobal 全局刷新
func RefreshGlobal() serializer.Response {
	err := RefreshDatabase()
	if err != nil {
		return serializer.CreateInternalErrorResponse("同步仓库失败", err)
	}
	return serializer.CreateSuccessResponse("", "同步仓库成功")
}

func RefreshDatabase() error {
	// 清空所有数据库
	if err := clearAll(); err != nil {
		w := fmt.Errorf("clearAll failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	// 遍历每个 Repo 配置项
	var repos []string
	for _, r := range conf.GlobalConfig.Repositories {
		repos = append(repos, filepath.Base(r.Root))
		err := loadLocalRepo(r)
		if err != nil {
			w := fmt.Errorf("loadLocalRepo failed: %w", err)
			utils.Log().Errorf("%v", w)
			return w
		}
	}
	// 授予管理员所有仓库的访问权限
	err := models.UpdateUserAllowedRepos(1, repos)
	if err != nil {
		w := fmt.Errorf("UpdateUserAllowedRepos to admin failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	return nil
}

// RefreshRepo 只刷新特定的Repo
func RefreshRepo(repoName string) serializer.Response {
	// 获取该 repo 的配置信息
	cfg := conf.GetDirectoryConfig(repoName)
	if cfg == nil {
		return serializer.CreateParamErrorResponse(errors.New("repo not exist"))
	}
	// 清空该 repo 相关的数据
	err := clearRepo(repoName)
	if err != nil {
		return serializer.CreateDBErrorResponse("", err)
	}
	err = loadLocalRepo(*cfg)
	if err != nil {
		return serializer.CreateInternalErrorResponse("loadLocalRepo failed", err)
	}
	utils.Log().Info("[Success] Refresh Local Repository")
	return serializer.CreateSuccessResponse("", "同步当前仓库成功")
}

// RefreshTag 只刷新特定的Tag
func RefreshTag(repoName string, tagName string) serializer.Response {
	// 获取该 repo 的配置信息
	cfg := conf.GetDirectoryConfig(repoName)
	if cfg == nil {
		return serializer.CreateParamErrorResponse(errors.New("repo not exist"))
	}
	// 清空 tag 下的 articles
	err := clearTag(repoName, tagName)
	if err != nil {
		return serializer.CreateDBErrorResponse("", err)
	}
	// 获取该目录的文章列表信息，重新生成文章
	tag, err := models.GetTag(repoName, tagName)
	if err != nil {
		return serializer.CreateDBErrorResponse("", err)
	}
	err = loadLocalTag(tag)
	if err != nil {
		return serializer.CreateInternalErrorResponse("loadLocalTag failed", err)
	}
	utils.Log().Info("[Success] Refresh Local Tag")
	return serializer.CreateSuccessResponse("", "同步当前标签成功")
}

// clearAll 清空所有 Article、Tag、Repo
func clearAll() error {
	// 清空 tag 数据库
	err := models.TruncateTags()
	if err != nil {
		w := fmt.Errorf("truncate tag db failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	// 清空 repo 数据库
	err = models.TruncateRepo()
	if err != nil {
		w := fmt.Errorf("truncate repo db failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	// 清空 article 数据库
	err = models.TruncateArticles()
	if err != nil {
		w := fmt.Errorf("truncate article db failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	return nil
}

// clearRepo 清空指定 Repo 及以下的 Tag 和 Article
func clearRepo(repoName string) error {
	// 清空指定 repo
	err := models.DeleteRepo(repoName)
	if err != nil {
		w := fmt.Errorf("delete repo failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	// 清空 repo 下的 tags
	err = models.DeleteTagsByRepo(repoName)
	if err != nil {
		w := fmt.Errorf("delete tags failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	// 清空 repo 下的 articles
	err = models.DeleteArticlesByRepo(repoName)
	if err != nil {
		w := fmt.Errorf("delete articles failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	return nil
}

// clearTag 清空指定 Tag 下的 Article
func clearTag(repoName string, tagName string) error {
	// 清空 repo 下的 articles
	err := models.DeleteArticlesByTag(repoName, tagName)
	if err != nil {
		w := fmt.Errorf("delete articles failed: %w", err)
		utils.Log().Errorf("%v", w)
		return w
	}
	return nil
}

// traverseRepo 遍历 Repo 目录，生成标签与文档的映射树
func traverseRepo(repo conf.DirectoryConfig) (tagPaths []string, tagPath2FileListMap utils.TagPath2FileInfo, err error) {
	tagPath2FileListMap = make(map[string][]*utils.FileInfo)
	var root = repo.Root
	var exclude = repo.Exclude

	// rootPrefix 用于拆分标签
	var rootPrefix = root
	if !strings.HasSuffix(rootPrefix, "/") {
		rootPrefix = rootPrefix + "/"
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
			if paths, ok := tagPath2FileListMap[dir]; ok {
				paths = append(paths, &utils.FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				tagPath2FileListMap[dir] = paths
			} else {
				paths = make([]*utils.FileInfo, 0, 10)
				paths = append(paths, &utils.FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				tagPath2FileListMap[dir] = paths
			}
		}
		return nil
	})
	return
}

func getArticlesInTag(tag models.Tag) (fileInfoList []*utils.FileInfo, err error) {
	var originalInfos []os.FileInfo
	originalInfos, err = ioutil.ReadDir(tag.Path)
	if err != nil {
		return
	}
	for _, info := range originalInfos {
		if strings.HasSuffix(info.Name(), ".md") {
			fileInfoList = append(fileInfoList, &utils.FileInfo{
				BriefName: strings.TrimSuffix(info.Name(), ".md"),
				Path:      filepath.Join(tag.Path, info.Name()),
			})
		}
	}
	return
}
