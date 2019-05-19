package Config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Mongo       Mongo       `json:"mongo,omitempty" bson:"mongo,omitempty"`
	StatusMicro StatusMicro `json:"statusMicro,omitempty" bson:"statusMicro,omitempty"`
}

type Mongo struct {
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Address string `json:"address,omitempty" bson:"address,omitempty"`
	Port    string `json:"port,omitempty" bson:"port,omitempty"`
}

type StatusMicro struct {
	Port string `json:"port,omitempty" bson:"port,omitempty"`
}

//Load configuration from file.
func GetAll() Config {
	jsonFile, err := os.Open(GetEnvRoute() + "Config/config." + getEnv() + ".json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	return config
}

func getEnv() string {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "dev"
	}
	return goEnv
}

func GetEnvRoute() string {
	goEnv := os.Getenv("GO_PROJECT_CONF_ROUTE")
	if goEnv == "" {
		goEnv = "./"
	}

	return goEnv
}
