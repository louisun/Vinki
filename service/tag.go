package service

import (
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
)

type TagArticleInfoView struct {
	models.TagView
	ArticleInfos []models.ArticleInfo
}

// GetRootTagInfosByRepo 获取 Repo 下一级 Tag 信息
func GetRootTagInfosByRepo(repoID uint64) serializer.Response {
	tags, err := models.GetRootTagInfosByRepo(repoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("当前仓库无标签可获取", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	return serializer.SuccessResponse(tags, "")
}

// GetTagViewByID 根据 TagID 获取 TagArticleInfoView
func GetTagViewByID(tagID uint64, flat bool) serializer.Response {
	var (
		tagView models.TagView
		err     error
	)
	if flat {
		tagView, err = models.GetFlatTagViewByID(tagID)
	} else {
		tagView, err = models.GetTagViewByID(tagID)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.ParamErrorResponse("tag id 不存在", err)
		}
		return serializer.DBErrorResponse("", err)
	}
	articles, err := models.GetArticleInfosByTagID(tagID)
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
