package models

type Repo struct {
	ID   uint64 `gorm:"primary_key"`
	Name string `gorm:"type:varchar(50);unique_idx"`
	Path string `gorm:"type:varchar(200);not null"` // 根目录路径
}

func (Repo) TableName() string {
	return "repo"
}

func AddRepo(repo *Repo) error {
	result := DB.Create(repo)
	return result.Error
}

func GetRepoNames() ([]string, error) {
	l := make([]string, 0, 10)
	result := DB.Model(&Repo{}).Pluck("name", &l)
	return l, result.Error
}

func DeleteRepo(repoName string) error {
	result := DB.Where("name = ?", repoName).Delete(&Repo{})
	return result.Error
}

func TruncateRepo() (err error) {
	result := DB.Delete(&Repo{})
	return result.Error
}
