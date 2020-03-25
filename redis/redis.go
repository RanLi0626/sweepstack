package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetConn() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 6379))
	if err != nil {
		log.Println("connect redis error", err)
		return nil, err
	}

	return conn, nil
}

//InitRedis used to init award_time and award_remain_num
func InitRedis() error {
	conn, err := GetConn()
	if err != nil {
		log.Println("conn is nil")
		return err
	}
	defer conn.Close()

	conn.Send("ZADD", "award_remain_num", 200, "A")
	conn.Send("ZADD", "award_remain_num", 400, "B")
	conn.Send("ZADD", "award_remain_num", 800, "C")
	conn.Send("HSET", "award_time", "A", time.Now().Format(time.RFC3339))
	conn.Send("HSET", "award_time", "B", time.Now().Format(time.RFC3339))
	conn.Send("HSET", "award_time", "C", time.Now().Format(time.RFC3339))
	conn.Flush()

	for i := 0; i < 3; i++ {
		_, err := conn.Receive()
		if err != nil {
			log.Printf("conn send error, %s", err)
			return err
		}
	}

	return nil
}
