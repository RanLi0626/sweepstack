package handler

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	red "sweepstake/redis"

	"github.com/gomodule/redigo/redis"
)

var (
	layout    string = "2006-01-02 15:04:05"
	startTime string = "2020-03-25 17:00:00"
	endTime   string = "2020-03-25 19:00:00"
)

type award struct {
	name             string
	remainedNum      int64
	lastReleasedTime time.Time
}

// Draw is the handler to draw the award
func Draw() {
	http.HandleFunc("/draw", draw)
	http.ListenAndServe(":8080", nil)
}

func draw(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	usernames, ok := params["username"]
	if !ok {
		log.Println("username is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var username string
	username = usernames[0]

	award, err := winCheck()
	if err != nil {
		log.Printf("nothing, err :%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("win nothing"))
		return
	}
	if award != nil {
		log.Println("win")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(username + ", you win"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("lose"))
	}

	return
}

func winCheck() (*award, error) {
	award, err := getRamdomAward()
	if err != nil {
		log.Printf("err in winCheck(), err : %v", err)
		return nil, err
	}
	log.Printf("award %v", award)

	end, err := time.Parse(layout, endTime)
	if err != nil {
		log.Printf("err in winCheck(), err : %v", err)
		return nil, err
	}
	start, err := time.Parse(layout, startTime)
	if err != nil {
		log.Printf("err in winCheck(), err : %v", err)
		return nil, err
	}

	deltaTime := end.Sub(start).Nanoseconds() / getTotalPrizeNum
	random := rand.New(rand.NewSource(end.Sub(award.lastReleasedTime).Nanoseconds()))

	nextReleasedTime := start.UnixNano() + deltaTime*getReleasedNum(*award) + int64(random.Int())%deltaTime
	log.Printf("nextReleasedTime %v, now %v", nextReleasedTime, time.Now().UnixNano())
	if time.Now().UnixNano() < nextReleasedTime {
		return nil, errors.New("failed")
	}

	conn, err := red.GetConn()
	if err != nil {
		log.Printf("err in winCheck(), err : %v", err)
		return nil, err
	}
	defer conn.Close()

	conn.Send("WATCH", "award_remain_num")
	conn.Send("MULTI")
	conn.Send("HSET", "award_time", award.name, time.Now().Unix())
	conn.Send("ZADD", "award_remain_num", award.name, award.remainedNum-1)
	conn.Send("EXEC")

	for i := 0; i < 3; i++ {
		_, err := conn.Receive()
		if err != nil {
			log.Printf("conn send error, %s", err)
			return nil, err
		}
	}

	return award, nil
}

func getRamdomAward() (*award, error) {
	conn, err := red.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Get remained num.
	result, err := redis.StringMap(conn.Do("ZRANGE", "award_remain_num", 0, -1, "WITHSCORES"))
	if err != nil {
		log.Printf("err in getRemainedNum() from redis, err : %v", err.Error())
		return nil, err
	}
	var totalRemainedNum int64
	for _, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Printf("err in getRemainedNum(), err : %v", err)
			return nil, err
		}
		totalRemainedNum = totalRemainedNum + remainedNum
	}
	log.Printf("totalRemaindNum : %v, result : %v", totalRemainedNum, result)

	// Get random award.
	random := rand.New(rand.NewSource(totalRemainedNum))
	num := random.Int63n(totalRemainedNum)
	log.Printf("num : %V", num)

	var a *award
	var total int64
	for k, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if remainedNum == 0 {
			continue
		}
		total = total + remainedNum
		if err != nil {
			log.Printf("err in getRamdomAward(), err : %v", err)
			return nil, err
		}
		if num-total < 0 {
			a = &award{name: k, remainedNum: remainedNum}
			break
		}
	}

	if a == nil {
		return nil, errors.New("err in getRemainedNum(), got nothing")
	}

	// Get lastUpdateTime.
	lastUpdateTimeStr, err := redis.String(conn.Do("HGET", "award_time", a.name))
	if err != nil {
		log.Printf("err in getLastUpdateTime(), err : %v", err)
		return nil, err
	}
	lastUpdateTime, err := time.Parse(time.RFC3339, lastUpdateTimeStr)
	if err != nil {
		log.Printf("err in getLastUpdateTime(), err : %v", err)
		return nil, err
	}
	a.lastReleasedTime = lastUpdateTime

	return a, nil
}

func getTotalPrizeNum() int64 {
	return 200 + 400 + 800
}

func getReleasedNum(a award) int64 {
	// TODO get from redis
	if a.name == "A" {
		return 200 - a.remainedNum
	}
	if a.name == "B" {
		return 400 - a.remainedNum
	}
	if a.name == "C" {
		return 800 - a.remainedNum
	}
	return 0
}
