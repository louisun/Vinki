package utils

import "fmt"

type FileInfo struct {
	BriefName string
	Path      string
}

func (fi *FileInfo) String() string {
	return fmt.Sprintf("[%s] -> [%s]", fi.BriefName, fi.Path)
}

var Tag2FilePath = make(FileMap)

type FileMap map[string][]*FileInfo

func (m FileMap) String() string {
	var s string
	for k, v := range m {
		s += fmt.Sprintf("tag: %s\n", k)
		for _, fi := range v {
			s += fmt.Sprintf("\t%s\n", fi)
		}
	}
	return s
}
