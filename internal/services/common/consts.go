package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var APIVersion string

type ConfigType struct {
	NatsURL string `envconfig:"URL" default:"127.0.0.1"`
	Port    int    `envconfig:"PORT" default:"8000"`

	GracefulStop        uint          `envconfig:"STOP_TIMEOUT" default:"10"`
	GracefulStopTimeout time.Duration `ignored:"true"`

	StartTimeout         uint          `envconfig:"START_TIMEOUT" default:"10"`
	GracefulStartTimeout time.Duration `ignored:"true"`

	APIShutdown    uint          `envconfig:"API_SHUTDOWN" default:"10"`
	APIShutdownDur time.Duration `ignored:"true"`

	BackendShutdown    uint          `envconfig:"BACKEND_SHUTDOWN" default:"10"`
	BackendShutdownDur time.Duration `ignored:"true"`

	HelloSubject string `ignored:"true"`

	QueueGroup string `envconfig:"QUEUE" default:"ru.example.queue"`
}

var Config ConfigType

func init() {
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}

	Config.GracefulStopTimeout = time.Duration(int64(time.Second) * int64(Config.GracefulStop))
	Config.GracefulStartTimeout = time.Duration(int64(time.Second) * int64(Config.StartTimeout))

	Config.APIShutdownDur = time.Duration(int64(time.Second) * int64(Config.APIShutdown))
	Config.BackendShutdownDur = time.Duration(int64(time.Second) * int64(Config.BackendShutdown))

	Config.HelloSubject = fmt.Sprintf("api.%v.hello", APIVersion)

	data, _ := json.Marshal(&Config)
	log.Printf("Configuration: %#v\n", string(data))
}
