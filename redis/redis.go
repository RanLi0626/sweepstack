package redis

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func getConn() redis.Conn {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 6379))
	if err != nil {
		log.Println("connect redis error", err)
		return nil
	}

	return conn
}

//InitRedis used to init award_time and award_remain_num
func InitRedis() error {
	fmt.Println("test")
	conn := getConn()
	if conn == nil {
		log.Println("conn is nil")
		return errors.New("conn is nil")
	}
	defer conn.Close()

	conn.Send("ZADD", "award_remain_num", 20, "A")
	conn.Send("ZADD", "award_remain_num", 40, "B")
	conn.Send("ZADD", "award_remain_num", 80, "C")
	conn.Send("HSET", "award_time", "A", time.Now())
	conn.Send("HSET", "award_time", "B", time.Now())
	conn.Send("HSET", "award_time", "C", time.Now())
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
