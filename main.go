package main

import (
	"sweepstake/conf"
	"sweepstake/handler"
	"sweepstake/redis"
)

func init() {
	conf.InitConf()
	redis.InitRedis()
}

func main() {
	handler.Draw()
}
