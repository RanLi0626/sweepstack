package handler

import (
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
	startTime string = "2020-03-20 00:00:00"
	endTime   string = "2020-03-21 00:00:00"
)

type award struct {
	name        string
	remainedNum int64
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

	winCheck()

	w.Write([]byte(username))
}

func getRamdomAward() (*award, error) {
	conn, err := red.GetConn()
	if err != nil {
		return nil, err
	}
	conn.Flush()

	// Get remained num.
	result, err := redis.StringMap(conn.Do("ZRANGE", "award_remain_num", 0, -1, "WITHSCORE"))
	if err != nil {
		return nil, err
	}
	var totalRemainedNum int64
	for _, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		totalRemainedNum = totalRemainedNum + remainedNum
	}

	// Get random award.
	random := rand.New(rand.NewSource(totalRemainedNum))
	num := random.Int63n(totalRemainedNum)

	for k, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		if num-remainedNum < 0 {
			return &award{name: k, remainedNum: remainedNum}, nil
		}
	}

	return nil, nil
}

func winCheck() error {
	award, err := getRamdomAward()
	if err != nil {
		return err
	}

	end, err := time.Parse(layout, endTime)
	if err != nil {
		return err
	}
	start, err := time.Parse(layout, startTime)
	if err != nil {
		return err
	}

	totalPrizeNum := getTotalPrizeNum()
	deltaTime := end.Sub(start).Nanoseconds() / totalPrizeNum
	random := rand.New(rand.NewSource(end.Sub(getLastReleasedTime()).Nanoseconds()))

	nextReleasedTime := start.UnixNano() + deltaTime*getReleasedNum() + int64(random.Int())
	if time.Now().UnixNano() >= nextReleasedTime {
		// TODO win the prize
	}

	conn, err := red.GetConn()
	if err != nil {
		return err
	}
	// TODO update redis remained num

	// TODO update redis released time

	return nil
}

func getTotalPrizeNum() int64 {
	return 10 + 20 + 30
}

func getReleasedNum() int64 {
	// TODO get from redis
	return 20
}

func getLastReleasedTime() time.Time {
	// TODO get from redis
	return time.Now()
}
