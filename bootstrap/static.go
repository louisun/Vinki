package bootstrap

import (
	"net/http"

	"github.com/louisun/vinki/pkg/utils"

	"github.com/gin-contrib/static"
	_ "github.com/louisun/vinki/statik"
	"github.com/rakyll/statik/fs"
)

// GinFS 实现了 ServeFileSystem 和 FileSystem 接口
type GinFS struct {
	FS http.FileSystem
}

var StaticFS static.ServeFileSystem

func (b *GinFS) Open(name string) (http.File, error) {
	return b.FS.Open(name)
}

func (b *GinFS) Exists(prefix string, filepath string) bool {
	if _, err := b.FS.Open(filepath); err != nil {
		return false
	}

	return true
}

// InitStatic 初始化静态资源的文件系统
func InitStatic() {
	var err error

	if utils.ExistsDir(utils.RelativePath("statics")) {
		utils.Log().Info("检测到 statics 目录存在，将使用此目录下的静态资源文件")
		// 使用 gin-contrib 的 static.LoadFile 创建静态文件系统
		StaticFS = static.LocalFile(utils.RelativePath("statics"), false)
	} else {
		// 不存在 static 目录，加载 statik 注册的数据（Zip压缩的字符串），创建文件系统
		StaticFS = &GinFS{}
		StaticFS.(*GinFS).FS, err = fs.New()
		if err != nil {
			utils.Log().Panicf("无法初始化静态资源, %s", err)
		}
	}
}
