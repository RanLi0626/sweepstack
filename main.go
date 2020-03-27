package main

import (
	"sweepstake/conf"
	"sweepstake/handler"
	"sweepstake/redis"
)

var Conf *conf.Yaml

func init() {
	conf.InitConf()
	redis.InitRedis()
}

func main() {
	handler.Draw()
}
