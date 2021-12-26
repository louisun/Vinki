package service

import (
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/model"
	"github.com/louisun/vinki/pkg/serializer"
)

type tagArticleInfoView struct {
	model.TagView
	ArticleInfos []string
}

// GetTopTagInfosByRepo 获取 Repo 的一级 Tag 信息
func GetTopTagInfosByRepo(repoName string) serializer.Response {
	tags, err := model.GetTopTagInfosByRepo(repoName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.CreateGeneralParamErrorResponse("当前仓库无标签可获取", err)
		}

		return serializer.CreateDBErrorResponse("", err)
	}

	return serializer.CreateSuccessResponse(tags, "")
}

// GetTagArticleView 根据 TagID 获取 TagArticleInfoView
func GetTagArticleView(repoName string, tagName string, flat bool) serializer.Response {
	var (
		tagView model.TagView
		err     error
	)

	if flat {
		tagView, err = model.GetFlatTagView(repoName, tagName)
	} else {
		tagView, err = model.GetTagView(repoName, tagName)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.CreateGeneralParamErrorResponse("tag 不存在", err)
		}

		return serializer.CreateDBErrorResponse("", err)
	}

	articles, err := model.GetArticleList(repoName, tagName)

	if err != nil {
		return serializer.CreateDBErrorResponse("", err)
	}

	sort.Sort(model.Articles(articles))

	view := tagArticleInfoView{
		TagView:      tagView,
		ArticleInfos: articles,
	}

	return serializer.CreateSuccessResponse(view, "")
}
