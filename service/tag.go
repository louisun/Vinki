package service

import (
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
)

type TagArticleInfoView struct {
	models.TagView
	ArticleInfos []string
}

// GetTopTagInfosByRepo 获取 Repo 的一级 Tag 信息
func GetTopTagInfosByRepo(repoName string) serializer.Response {
	tags, err := models.GetTopTagInfosByRepo(repoName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("当前仓库无标签可获取", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	return serializer.SuccessResponse(tags, "")
}

// GetTagArticleView 根据 TagID 获取 TagArticleInfoView
func GetTagArticleView(repoName string, tagName string, flat bool) serializer.Response {
	var (
		tagView models.TagView
		err     error
	)
	if flat {
		tagView, err = models.GetFlatTagView(repoName, tagName)
	} else {
		tagView, err = models.GetTagView(repoName, tagName)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("tag 不存在", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	articles, err := models.GetArticleList(repoName, tagName)
	if err != nil {
		return serializer.DBErrorResponse("", err)
	}
	sort.Sort(models.Articles(articles))
	view := TagArticleInfoView{
		TagView:      tagView,
		ArticleInfos: articles,
	}
	return serializer.SuccessResponse(view, "")
}

// 添加标签名
func AddTags(tags []*models.Tag) error {
	err := models.AddTags(tags)
	return err
}

// truncateTags 清空 Tag
func truncateTags() error {
	err := models.TruncateTags()
	return err
}

// deleteTagsByRepo 清空某个 repo 下的 Tags
func deleteTagsByRepo(repoID uint64) error {
	err := models.DeleteTagsByRepo(repoID)
	return err
}
