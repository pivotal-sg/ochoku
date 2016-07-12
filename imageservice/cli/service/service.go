package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-platform/config"
	"github.com/micro/go-platform/config/source/file"
	"github.com/pivotal-sg/ochoku/imageservice"
	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
)

var (
	configFile  = filepath.Join(os.TempDir(), "imageservice.config")
	ServiceName = "pivotal.io.ochoku.imageservice"
	Version     = "0.1.0"
)

func removeFile() error {
	return os.Remove(configFile)
}

func main() {
	var conf string = `{
		"aws": {
			"id": "asdf",
			"secret": "secret-asdf",
			"token": "token-asdf",
			"region": "ap-southeast-1",
			"bucket": "ochoku"
		}
	}`

	if err := ioutil.WriteFile(configFile, []byte(conf), 0600); err != nil {
		log.Println("Can't write the conf file")
		panic(err)
	}
	defer removeFile()

	// Create a config instance
	config := config.NewConfig(
		// Poll every minute for changes
		config.PollInterval(time.Minute),
		// Use file as a config source
		// Multiple sources can be specified
		config.WithSource(file.NewSource(config.SourceName(configFile))),
	)

	service := micro.NewService(
		micro.Name(ServiceName),
		micro.Version(Version),
	)

	service.Init()

	s3store := imageservice.S3Store{
		Config: config,
	}

	is := &imageservice.ImageService{
		Config:    config,
		DataStore: make([]proto.ImageList, 0, 0),
		FileStore: s3store,
	}

	proto.RegisterImageStorerHandler(service.Server(), is)
	service.Run()
}
