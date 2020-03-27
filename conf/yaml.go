package conf

// Yaml struct of yaml
type Yaml struct {
	RedisConf    Redis    `yaml:"redis"`
	InitTimeConf InitTime `yaml:"initTime"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type InitTime struct {
	StartTime string `yaml:"starttime"`
	EndTime   string `yaml:"endtime"`
}
