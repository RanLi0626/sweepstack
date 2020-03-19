package handler

import (
	"log"
	"net/http"
)

//Draw is the handler to draw the award
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
	w.Write([]byte("test"))
}
