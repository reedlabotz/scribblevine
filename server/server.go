package scribblevine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"appengine"
	"appengine/datastore"
)

type Move struct {
	UserId string
	Type int
	Data []byte
	Timestamp time.Time
}

type Game struct {
	UserId string
	Started time.Time
  Players []string
}

type GamePlayer struct {
	UserId string
}

func init() {
	http.HandleFunc("/api/alive", alive)
	http.HandleFunc("/api/game/start", gameStart)
	http.HandleFunc("/api/game/get", gameGet)
	http.HandleFunc("/api/game/move", gameMove)
}

func alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "yes")
}

type GameStartResponse struct {
	Success bool
	Id string
}

func gameStart(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	userId := r.FormValue("UserId")
	players := strings.Split(r.FormValue("Players"),",")
	
	g := Game{
		UserId: userId,
		Started: time.Now(),
		Players: players,
	}
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "game", nil), &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GameStartResponse{
		Success: true,
		Id: key.Encode(),
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

type GameGetResponse struct {
	Type int
	Data string
}

func gameGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := r.FormValue("Id")
	gameKey, err := datastore.DecodeKey(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := datastore.NewQuery("move").Ancestor(gameKey).Order("-Timestamp").Limit(1)
	moves := make([]Move, 0, 1)
	if _, err := q.GetAll(c, &moves); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := GameGetResponse{
		Type: moves[0].Type,
		Data: string(moves[0].Data),
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

type GameMoveResponse struct {
	Success bool
	Id string
}

func gameMove(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := r.FormValue("GameId")

	gameKey, err := datastore.DecodeKey(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	moveType,_ := strconv.Atoi(r.FormValue("Type"))
	moveData := []byte(r.FormValue("Data"))

	userId := r.FormValue("UserId")
	m := Move{
		UserId: userId,
		Timestamp: time.Now(),
		Type: moveType,
		Data: moveData,
	}	

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "move", gameKey), &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GameMoveResponse{
		Success: true,
		Id: key.Encode(),
	}
	
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}