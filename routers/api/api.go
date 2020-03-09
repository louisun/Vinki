package api

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/vinki/service"

	"github.com/gin-gonic/gin"
	"github.com/vinki/utils"
)

// GET 主页
func GetHomePage(c *gin.Context) {
	article, err := service.GetArticle("Home", "Vinki")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Redirect(http.StatusTemporaryRedirect, "/refresh")
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(article.HtmlContent))
}

// GET 标签主页
func GetTagPage(c *gin.Context) {
	tag, err := service.GetTagHtml(c.Param("tag"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(tag.HtmlContent))
}

// GET 文章页
func GetWikiPage(c *gin.Context) {
	article, err := service.GetArticle(c.Param("tag"), c.Param("title"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(article.HtmlContent))
}

// GET 刷新
func Refresh(c *gin.Context) {
	err := service.Refresh()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/wiki")
}

// 加载 Markdown -> 渲染 Html 返回
func loadPage(c *gin.Context, mdPath string, number int, currentTag string, files []string) {
	tags, err := service.GetAllTags()
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
