package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Conf struct of yaml
type Conf struct {
	RedisConf    Redis    `yaml:"redis"`
	InitTimeConf InitTime `yaml:"initTime"`
	AwardConf    Award    `yaml:"award"`
}

// Redis is the config for redis
type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// InitTime is the config for sweepstake start and end time
type InitTime struct {
	StartTime string `yaml:"starttime"`
	EndTime   string `yaml:"endtime"`
	Layout    string `yaml:"layout"`
}

// Award is the config for prize
type Award struct {
	A int64 `yaml:"A"`
	B int64 `yaml:"B"`
	C int64 `yaml:"C"`
}

var (
	// RedisConf is the redis conf from yaml
	RedisConf Redis
	// InitTimeConf is the init time conf from yaml
	InitTimeConf InitTime
	// AwardConf is the award info conf from yaml
	AwardConf Award
)

// InitConf used to init config yaml file
func InitConf() {
	var config *Conf

	yamlFile, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}

	log.Println("conf:", config)

	RedisConf = config.RedisConf
	InitTimeConf = config.InitTimeConf
	AwardConf = config.AwardConf
}
