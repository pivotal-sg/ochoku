package reviews

import (
	"time"

	"github.com/micro/go-platform/config"
	"github.com/micro/go-platform/config/source/file"
)

const DefaultStorageFileName = "./reviews.db"
const DefaultConfigFileName = "./reviews.config"

var ReviewConfig config.Config
var ConfigFileName string

func Config(fileName string) config.Config {
	return config.NewConfig(
		// Poll every minute for changes
		config.PollInterval(time.Minute),
		// Use file as a config source
		// Multiple sources can be specified
		config.WithSource(file.NewSource(config.SourceName(fileName))),
	)
}
