package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vinki/templates"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var log = DefaultLog

// GenerateWikiHtml 从本地读取 Markdown，填充模板，生成完整的 Wiki Html bytes
func GenerateWikiHtml(mdPath string, number int, currentTag string, tags []string, files []string) (string, error) {
	html, err := markdownFile2Html(mdPath)
	if err != nil {
		return "", err
	}
	return templates.RenderWiki(string(html), number, currentTag, tags, files), nil
}

// MarkdownFile2Html 读取本地 Markdown 文件并渲染为 Html bytes
func markdownFile2Html(mdPath string) ([]byte, error) {
	if !strings.HasSuffix(mdPath, ".md") {
		return nil, errors.New("File suffix is not \".md\"")
	}
	mdFile, err := os.Open(mdPath)
	if err != nil {
		log.Errorf("Open mdFile %s failed: %v", mdPath, err)
		return nil, err
	}
	md, err := ioutil.ReadAll(mdFile)
	if err != nil {
		log.Errorf("Read mdFile %s failed", mdPath)
		return nil, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(md, parser, nil)
	return html, nil
}

// 遍历目录
func Traverse() error {
	// 清空FileMap
	Tag2FilePath = make(FileMap)
	var root = Config.Directory.Root
	var exclude = Config.Directory.Exclude
	var fold = Config.Directory.Fold
	log.Info("Start Load Local Files...")
	log.Infof("Root: %s\tExclude: %v\tFold: %v", root, exclude, fold)
	if !existsDir(root) {
		return errors.New("root directory does not exist")
	}

	var rootPrefix = root
	if !strings.HasSuffix(Config.Directory.Root, "/") {
		rootPrefix = Config.Directory.Root + "/"
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//跳过指定的目录
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || isDirectoryInList(info.Name(), exclude) {
				return filepath.SkipDir
			}
			return nil
		}

		// 遍历 Markdown 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			p := strings.TrimPrefix(path, rootPrefix)
			p = strings.TrimSuffix(p, ".md")
			l := strings.Split(p, "/")
			l = l[:len(l)-1]
			var tags = make([]string, 0, len(l))
			for _, t := range l {
				if !isDirectoryInList(t, fold) {
					tags = append(tags, t)
				}
			}
			tagName := strings.Join(tags, "-")
			if paths, ok := Tag2FilePath[tagName]; ok {
				paths = append(paths, &FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				Tag2FilePath[tagName] = paths
			} else {
				paths = make([]*FileInfo, 0, 10)
				paths = append(paths, &FileInfo{
					BriefName: strings.TrimSuffix(info.Name(), ".md"),
					Path:      path,
				})
				Tag2FilePath[tagName] = paths
			}
		}
		return nil
	})
	return nil
}

// 判断所给路径是否存在
func existsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断路径是否是目录
func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// 判断路径是否是文件
func isFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// 判断是否存在该目录
func existsDir(path string) bool {
	return existsPath(path) && isDir(path)
}

// 判断是否存在该文件
func existsFile(path string) bool {
	return existsPath(path) && isFile(path)
}

func getTags(path string) string {
	return ""

}
func isDirectoryInList(directory string, list []string) bool {
	for _, t := range list {
		if directory == t {
			return true
		}
	}
	return false
}
