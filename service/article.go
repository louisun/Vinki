package service

import (
	"github.com/vinki/db"
)

// 获取文章
func GetArticle(tag string, title string) (*db.Article, error) {
	article, err := db.GetArticleByTagAndTitle(tag, title)
	if err != nil {
		log.Errorf("Get article by tag and title failed: %v, tag: %s, title: %s", err, tag, title)
		return nil, err
	}
	return article, nil
}

// 获取同一 Tag 下的文章名列表
func GetArticleListByTag(tag string) ([]string, error) {
	articleNames, err := db.GetArticleListByTag(tag)
	if err != nil {
		log.Errorf("Get article list by tag failed: %v, tag: %s", err, tag)
		return nil, err
	}
	return articleNames, nil
}

// 添加 Articles
func AddArticles(articles []*db.Article) error {
	var err error
	for _, article := range articles {
		err = db.AddArticle(article)
		if err != nil {
			log.Errorf("Add articles failed: %v, article: %v", article, err)
			return err
		}
	}
	return nil
}
