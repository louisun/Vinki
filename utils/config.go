package utils

var Config = struct {
	Server struct {
		Port int `default:"8080"`
	}
	Database struct {
		Host     string `default:"vinki-postgres"`
		Port     int    `default:"5432"`
		User     string `default:"postgres"`
		Password string `default:"root"`
		Database string `default:"vinki"`
	}
	Directory struct {
		Root    string   `default:"/vinki/repository"`
		Exclude []string `required:"false"`
		Fold    []string `required:"false"`
	}
	Custom struct {
		Tag  string `default:"./conf/tag.md"`
		Home string `default:"./conf/home.md"`
	}
}{}
