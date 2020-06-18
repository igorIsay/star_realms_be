package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var addr = flag.String("addr", ":8080", "http service address")
var hubs = make(map[string]*Hub)

func route(hubs map[string]*Hub, w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "null")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Headers", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")

	var hubsPattern = regexp.MustCompile("^/hubs/(\\w+)\\?player=([1,2])")
	type HubData struct {
		Name string `json:"name"`
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "home.html")
		return
	}
	if r.URL.Path == "/hubs" && r.Method == "GET" {
		hubsList := make([]HubData, 0, len(hubs))
		for name := range hubs {
			hubsList = append(hubsList, HubData{Name: name})
		}
		var result []byte
		result, err := json.Marshal(hubsList)
		if err != nil {
			log.Println(err)
		}
		w.Write(result)
		return
	}
	if r.URL.Path == "/hubs" && r.Method == "POST" {
		var hubData HubData
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &hubData)
		if err != nil {
			log.Println(err)
		}
		_, ok := hubs[hubData.Name]
		if ok {
			http.Error(w, "Hub with such name already exists", http.StatusBadRequest)
			return
		}
		hub := newHub()
		hubs[hubData.Name] = hub
		w.WriteHeader(http.StatusOK)
		go hub.run()
		return
	}
	if hubsPattern.MatchString(r.URL.Path + "?" + r.URL.RawQuery) {
		matches := hubsPattern.FindStringSubmatch(r.URL.Path + "?" + r.URL.RawQuery)
		hub, ok := hubs[matches[1]]
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if matches[2] == "1" {
			serveWs(hub, FirstPlayer, w, r)
		}
		if matches[2] == "2" {
			serveWs(hub, SecondPlayer, w, r)
		}
		return
	}
	if r.Method == "OPTIONS" {
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func serveWs(hub *Hub, player PlayerId, w http.ResponseWriter, r *http.Request) {
	//TODO delete CheckOrigin reasigning
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, playerId: player, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route(hubs, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
