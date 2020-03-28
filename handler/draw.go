package handler

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	"sweepstake/conf"
	red "sweepstake/redis"
)

type award struct {
	name             string
	remainedNum      int64
	totalRemainedNum int64
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("got error"))
		log.Printf("got err %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if award == nil {
		w.Write([]byte("win nothing"))
		return
	}

	w.Write([]byte(username + ", you win"))
	return
}

func winCheck() (*award, error) {
	award, err := getRamdomAward()
	if err != nil {
		return nil, err
	}

	nextReleasedTime, err := getNextReleasedTime(award)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() < nextReleasedTime {
		return nil, errors.New("current time %v is before nextReleasedTime")
	}

	// Update redis
	conn, err := red.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.Send("WATCH", "award_remain_num")
	conn.Send("MULTI")
	conn.Send("HSET", "award_time", award.name, time.Now().Format(time.RFC3339))
	conn.Send("ZADD", "award_remain_num", award.name, award.remainedNum-1)
	conn.Send("EXEC")

	err = conn.Flush()
	if err != nil {
		return nil, err
	}

	return award, nil
}

func getRamdomAward() (*award, error) {
	conn, err := red.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	result, err := redis.StringMap(conn.Do("ZRANGE", "award_remain_num", 0, -1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	// Get remainedNum of all the awards.
	totalRemainedNum, err := getRemainedNum(result)
	if err != nil {
		return nil, err
	}

	// Get random award.
	random := rand.New(rand.NewSource(totalRemainedNum))
	randomNum := random.Int63n(totalRemainedNum)
	a, err := getAwardInfo(result, randomNum)
	if err != nil {
		return nil, err
	}

	// Get lastUpdateTime.
	lastUpdateTimeStr, err := redis.String(conn.Do("HGET", "award_time", a.name))
	if err != nil {
		return nil, err
	}
	lastUpdateTime, err := time.Parse(time.RFC3339, lastUpdateTimeStr)
	if err != nil {
		return nil, err
	}

	a.lastReleasedTime = lastUpdateTime
	a.totalRemainedNum = totalRemainedNum
	return a, nil
}

func getRemainedNum(result map[string]string) (int64, error) {
	var totalRemainedNum int64
	for _, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		totalRemainedNum = totalRemainedNum + remainedNum
	}

	if totalRemainedNum == 0 {
		return 0, errors.New("the awards are over")
	}

	return totalRemainedNum, nil
}

func getAwardInfo(result map[string]string, randomNum int64) (*award, error) {
	var a *award
	var total int64
	for k, v := range result {
		remainedNum, err := strconv.ParseInt(v, 10, 64)
		if remainedNum == 0 {
			continue
		}
		total = total + remainedNum
		if err != nil {
			return nil, err
		}
		if randomNum-total < 0 {
			a = &award{name: k, remainedNum: remainedNum}
			break
		}
	}

	if a == nil {
		return nil, errors.New("err in getRemainedNum(), got nothing")
	}

	return a, nil
}

func getNextReleasedTime(award *award) (int64, error) {
	end, err := time.ParseInLocation(conf.Conf.InitTimeConf.Layout, conf.Conf.InitTimeConf.EndTime, time.Local)
	if err != nil {
		return 0, err
	}
	start, err := time.ParseInLocation(conf.Conf.InitTimeConf.Layout, conf.Conf.InitTimeConf.StartTime, time.Local)
	if err != nil {
		return 0, err
	}

	e := end.Unix()
	s := start.Unix()

	deltaTime := (e - s) / getTotalAwardNum()
	random := rand.New(rand.NewSource(e - award.lastReleasedTime.Unix()))

	nextReleasedTime := s + deltaTime*(getTotalAwardNum()-award.totalRemainedNum) + int64(random.Int())%deltaTime

	return nextReleasedTime, nil
}

func getTotalAwardNum() int64 {
	aNum := conf.Conf.AwardConf.A
	bNum := conf.Conf.AwardConf.B
	cNum := conf.Conf.AwardConf.C

	return aNum + bNum + cNum
}
