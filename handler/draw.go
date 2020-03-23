package handler

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	layout    string = "2006-01-02 15:04:05"
	startTime string = "2020-03-20 00:00:00"
	endTime   string = "2020-03-21 00:00:00"
)

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

	// TODO getRamdomAward()
	// TODO winCheck()

	w.Write([]byte(username))
}

func getRamdomAward() {

}

func winCheck() {
	end, err := time.Parse(layout, endTime)
	if err != nil {
		return
	}
	start, err := time.Parse(layout, startTime)
	if err != nil {
		return
	}

	totalPrizeNum := getTotalPrizeNum()
	deltaTime := end.Sub(start).Nanoseconds() / totalPrizeNum
	random := rand.New(rand.NewSource(end.Sub(getLastReleasedTime()).Nanoseconds()))

	nextReleasedTime := start.UnixNano() + deltaTime*getReleasedNum() + int64(random.Int())
	if time.Now().UnixNano() >= nextReleasedTime {
		// TODO win the prize
	}
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
