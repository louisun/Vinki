package service

import (
	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
)

func GetAllRepoInfos() serializer.Response {
	repos, err := models.GetAllRepoNames()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.CreateGeneralParamErrorResponse("当前没有仓库可获取", err)
		}
		return serializer.CreateDBErrorResponse("", err)
	}
	return serializer.CreateSuccessResponse(repos, "")
}

// addRepo 添加 Repo
func addRepo(repo *models.Repo) error {
	err := models.AddRepo(repo)
	return err
}

func deleteRepo(repoName string) error {
	err := models.DeleteRepo(repoName)
	return err
}

func truncateRepo() error {
	err := models.TruncateRepo()
	return err
}
