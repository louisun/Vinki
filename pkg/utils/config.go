package utils

var Config = struct {
	Server struct {
		Port int `default:8080`
	}
	Database struct {
		Host     string `default:"127.0.0.1"`
		Port     int    `required:"true"`
		User     string `default:"postgres"`
		Password string `required:"true"`
		Database string `default:"vinki"`
	}
	Directory struct {
		Root    string   `required:"true"`
		Exclude []string `required:"false"`
		Fold    []string `required:"false"`
	}
}{}
