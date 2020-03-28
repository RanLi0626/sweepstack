package redis

import (
	"fmt"
	"log"
	"sweepstake/conf"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetConn() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Conf.RedisConf.Host, conf.Conf.RedisConf.Port))
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

	startTime, _ := time.Parse(conf.Conf.InitTimeConf.Layout, conf.Conf.InitTimeConf.StartTime)
	conn.Send("ZADD", "award_remain_num", 200, "A")
	conn.Send("ZADD", "award_remain_num", 400, "B")
	conn.Send("ZADD", "award_remain_num", 800, "C")
	conn.Send("HSET", "award_time", "A", startTime.Format(time.RFC3339))
	conn.Send("HSET", "award_time", "B", startTime.Format(time.RFC3339))
	conn.Send("HSET", "award_time", "C", startTime.Format(time.RFC3339))
	conn.Flush()

	for i := 0; i < 6; i++ {
		_, err := conn.Receive()
		if err != nil {
			log.Printf("conn send error, %s", err)
			return err
		}
	}

	return nil
}
