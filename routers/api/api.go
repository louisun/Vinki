package api

import (
	"net/http"

	"github.com/vinki/services"

	"github.com/gin-gonic/gin"
	"github.com/vinki/pkg/utils"
)

// GET 测试页
func GetTestPage(c *gin.Context) {
	loadPage(c, "./conf/test.md", 1, "测试", []string{"1. 测试", "2. 测试", "3. 测试"})
}

// GET 主页
func GetHomePage(c *gin.Context) {
	loadPage(c, "./conf/home.md", 1, "Home", nil)
}

// GET 标签主页
func GetTagPage(c *gin.Context) {
	files, err := services.GetArticleListByTag(c.Param("tag"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	loadPage(c, "./conf/tag.md", len(files), c.Param("tag"), files)
}

// GET 文章页
func GetWikiPage(c *gin.Context) {
	article, err := services.GetArticle(c.Param("tag"), c.Param("title"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(article.HtmlContent))
}

// POST 刷新
func Refresh(c *gin.Context) {
	err := services.Refresh()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/wiki")
}

// 加载 Markdown -> 渲染 Html 返回
func loadPage(c *gin.Context, mdPath string, number int, currentTag string, files []string) {
	tags, err := services.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "内部错误")
		return
	}
	html, err := utils.GenerateWikiHtml(mdPath, number, currentTag, tags, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "内部错误")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
