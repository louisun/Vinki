package models

type Repo struct {
	ID   uint64 `gorm:"primary_key"`
	Name string `gorm:"type:varchar(50);unique_idx"`
	Path string `gorm:"type:varchar(200);not null"` // 根目录路径
}

func (Repo) TableName() string {
	return "repo"
}

type RepoInfo struct {
	ID   uint64
	Name string
}

func GetRepos() ([]Repo, error) {
	l := make([]Repo, 0, 10)
	result := DB.Model(&Repo{}).Find(&l)
	return l, result.Error
}

func GetRepoInfos() ([]RepoInfo, error) {
	l := make([]RepoInfo, 0, 10)
	result := DB.Model(&Repo{}).Select("id, name").Scan(&l)
	return l, result.Error
}

func GetRepoByID(repoID uint64) (repo Repo, err error) {
	result := DB.First(&repo, repoID)
	return repo, result.Error
}

func DeleteRepoByID(repoID uint64) error {
	result := DB.Where("id = ?", repoID).Delete(&Repo{})
	return result.Error
}

func TruncateRepo() (err error) {
	result := DB.Delete(&Repo{})
	return result.Error
}

func AddRepo(repo *Repo) error {
	result := DB.Create(repo)
	return result.Error
}
