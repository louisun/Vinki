package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// FileInfo 文件信息
type FileInfo struct {
	BriefName string
	Path      string
}

func (fi *FileInfo) String() string {
	return fmt.Sprintf("[%s] -> [%s]", fi.BriefName, fi.Path)
}

// Tag path 到 FileInfo 的映射，用于创建文档
type TagPath2FileInfo map[string][]*FileInfo

func (m TagPath2FileInfo) String() string {
	var s string
	for k, v := range m {
		s += fmt.Sprintf("tag path: %s\n", k)
		for _, fi := range v {
			s += fmt.Sprintf("\t%s\n", fi)
		}
	}
	return s
}

// MarkdownFile2Html 读取本地 Markdown 文件并渲染为 Html bytes
func RenderMarkdown(mdPath string) ([]byte, error) {
	if !strings.HasSuffix(mdPath, ".md") {
		return nil, errors.New("File suffix is not \".md\"")
	}
	mdFile, err := os.Open(mdPath)
	if err != nil {
		Log().Errorf("Open mdFile %s failed: %v", mdPath, err)
		return nil, err
	}
	md, err := ioutil.ReadAll(mdFile)
	if err != nil {
		Log().Errorf("Read mdFile %s failed", mdPath)
		return nil, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(md, p, nil)
	return html, nil
}

// ExistsPath 判断所给路径是否存在
func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断路径是否是目录
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// IsFile 判断路径是否是文件
func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// ExistsDir 判断是否存在该目录
func ExistsDir(path string) bool {
	return ExistsPath(path) && IsDir(path)
}

// ExistsFile 判断是否存在该文件
func ExistsFile(path string) bool {
	return ExistsPath(path) && IsFile(path)
}

// IsDirectoryInList 判断目录是否在列表中
func IsDirectoryInList(directory string, list []string) bool {
	for _, t := range list {
		if directory == t {
			return true
		}
	}
	return false
}

// 获取运行环境的相对路径
func RelativePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	e, _ := os.Executable()
	return filepath.Join(filepath.Dir(e), name)
}
