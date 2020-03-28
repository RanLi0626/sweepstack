package conf

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Yaml struct of yaml
type Yaml struct {
	RedisConf    Redis    `yaml:"redis"`
	InitTimeConf InitTime `yaml:"initTime"`
	AwardConf    Award    `yaml:"award"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type InitTime struct {
	StartTime string `yaml:"starttime"`
	EndTime   string `yaml:"endtime"`
	Layout    string `yaml:"layout"`
}

type Award struct {
	A int64 `yaml:"A"`
	B int64 `yaml:"B"`
	C int64 `yaml:"C"`
}

var Conf *Yaml

func InitConf() {
	yamlFile, err := ioutil.ReadFile("test.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	log.Println("conf:", Conf)
}
