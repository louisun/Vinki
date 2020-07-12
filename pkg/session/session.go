package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/pkg/utils"
)

// GetSession 从 session 获取 key 对应的值
func GetSession(c *gin.Context, key string) interface{} {
	s := sessions.Default(c)
	return s.Get(c)
}

// SetSession 保存键值对到 session 中
func SetSession(c *gin.Context, kvMap map[string]interface{}) {
	s := sessions.Default(c)

	for key, value := range kvMap {
		s.Set(key, value)
	}

	err := s.Save()
	if err != nil {
		utils.Log().Warningf("无法设置 session: %s", err)
	}
}

// DeleteSession 删除 session
func DeleteSession(c *gin.Context, key string) {
	s := sessions.Default(c)
	s.Delete(key)
	err := s.Save()

	if err != nil {
		utils.Log().Warningf("无法删除 session key: %s, err: %s", key, err)
	}
}

// ClearSession 清空 session
func ClearSession(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	err := s.Save()

	if err != nil {
		utils.Log().Warningf("无法清空 session: %s", err)
	}
}
