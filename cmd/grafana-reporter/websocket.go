package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	t "time"

	"github.com/gorilla/websocket"
)

type subscription struct {
	conn     *websocket.Conn
	apitoken string
	params   url.Values
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var store = map[*websocket.Conn]int{}

var subscriptionChan = make(chan subscription)

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	apiToken := r.Header.Get("apitoken")
	params := r.URL.Query()

	//upgrade request to ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Cannot upgrade request to ws, , error: %+v", err)
		return
	}

	store[ws] = 0

	subscriptionChan <- subscription{ws, apiToken, params}
	for {
		t.Sleep(t.Duration(10) * t.Second)
		_, ok := store[ws]
		if !ok {
			break
		}
	}

}

// KeepAlive is for keeping alive websocket connection by sending ping pong
func KeepAlive(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				closeErr := err.(*websocket.CloseError)
				log.Printf("Socket closed %v; stopping reader\n", *closeErr)
			default:
				log.Println("Error in reader:", err)
			}
			return
		}
		if string(message) == "ping" {
			if err := ws.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				log.Printf("error writing pong to ws, error %+v: ", err)
			}
		}
	}
}

//reportHub runs in separate go routine and calls generate report
func reportHub() {
	for {
		select {
		case sub := <-subscriptionChan:
			go KeepAlive(sub.conn)
			go func(s subscription) {
				defer func() {
					delete(store, s.conn)
				}()
				apiToken := s.apitoken
				params := s.params
				dashboardUID := params.Get("dashboard")
				httpc := http.Client{
					Timeout: 10 * t.Minute,
				}

				url, err := url.Parse("http://localhost:8686/api/v5/report/" + dashboardUID)
				url.RawQuery = params.Encode()
				urlStr := url.String()

				req, err := http.NewRequest("GET", urlStr, nil)
				req.Header.Add("apitoken", apiToken)
				if err != nil {
					return
				}
				resp, err := httpc.Do(req)
				if err != nil {
					log.Println("Request to generate report failed , error: %+v", err)
					return
				}
				defer resp.Body.Close()

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Failed to read response body, error: %+v", err)
					return
				}
				if err := sub.conn.WriteMessage(websocket.BinaryMessage, bodyBytes); err != nil {
					log.Printf("Error writing to ws, error %+v: ", err)
				} else {
					log.Printf("Websocket write successful. Closing connection")
				}
				t.Sleep(t.Duration(5) * t.Second)
				sub.conn.Close()
			}(sub)

		}
	}
}
