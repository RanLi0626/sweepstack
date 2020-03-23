package main

import (
	"sweepstack/handler"
	"sweepstack/redis"
)

func init() {
	redis.InitRedis()
}

func main() {
	handler.Draw()
}
