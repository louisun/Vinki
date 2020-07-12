package conf

import (
	"path/filepath"

	"github.com/jinzhu/configor"
	"github.com/louisun/vinki/pkg/utils"
)

// GlobalConfig 全局配置变量
var GlobalConfig = struct {
	Database     DatabaseConfig
	Redis        RedisConfig
	System       SystemConfig
	Repositories []DirectoryConfig
}{}

// 数据库配置
type DatabaseConfig struct {
	Type     string `default:"UNSET"` // 数据库类型
	User     string
	Password string
	Host     string
	Port     uint
	DBName   string `default:"vinki"`    // 数据库名
	DBFile   string `default:"vinki.db"` // SQLite 数据库文件名
}

// Redis 配置
type RedisConfig struct {
	Host string
	Port uint
	DB   string `default:"0"`
}

// 系统配置
type SystemConfig struct {
	Debug         bool   `default:"false"` // 调试模式
	Port          uint   `default:"6166"`  // 监听端口
	SessionSecret string `default:"session-vinki-2020"`
	HashIDSalt    string `default:"hash-salt-2020"`
}

// 目录
type DirectoryConfig struct {
	Root    string   `default:"/vinki/repository"` // 根目录路径
	Exclude []string `required:"false"`            // 排除的文件、目录列表
	Fold    []string `required:"false"`            // 折叠的目录列表
}

// Init 从配置文件中初始化配置
func Init(path string) {
	if !utils.ExistsFile(path) {
		utils.Log().Panicf("配置文件路径不存在: %s", path)
	}

	err := configor.Load(&GlobalConfig, path)

	if err != nil {
		utils.Log().Panicf("加载配置文件失败, %s", err)
	}
}

// GetDirectoryConfig 获取相应仓库配置
func GetDirectoryConfig(repoName string) *DirectoryConfig {
	for _, repo := range GlobalConfig.Repositories {
		if filepath.Base(repo.Root) == repoName {
			return &repo
		}
	}

	return nil
}
