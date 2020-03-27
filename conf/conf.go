package conf

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func InitConf() *Yaml {
	conf := new(Yaml)
	yamlFile, err := ioutil.ReadFile("test.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
	log.Println("conf:", conf)

	return conf
}
