package main

import (
	"goexcersise/sweepstack/handler"
	"goexcersise/sweepstack/redis"
)

func init() {
	redis.InitRedis()
}

func main() {
	handler.Draw()
}
