package conf

import (
	"testing"

	"github.com/louisun/vinki/pkg/utils"
)

func TestInitConfig(t *testing.T) {
	Init("/Users/louisun/go/src/github.com/louisun/vinki/conf/config.yml")
	utils.PrettyPrint(&GlobalConfig)
}
