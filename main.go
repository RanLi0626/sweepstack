package main

import (
	"sweepstake/handler"
	"sweepstake/redis"
)

func init() {
	redis.InitRedis()
}

func main() {
	handler.Draw()
}
