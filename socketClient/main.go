package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:5001", "http service address")
var auth = flag.String("Authorization", "Bearer Zy0X31QIqIxGN2R3m5wONc8XeeAc-8MsEaaB5KQT_Ovx2KrZJygmjjsLS6D8MKJoSX60UaT879ftFjYrbSw2pc0GfIE0fD4JZwzno-3vBQwyy4xrYp046ZEFsI0OoboUc5XH8Furml8dz-TKtda3YckOd5ftZ2EWAhgtK2UHX5J-5jyN6mwEF0Ocdm5lHPrKdbmr3KcG_T5zOGE1t8Uw8yIHBltdhQubzCfvr5dKeZOcjLqjOG3meJUQe2S5cCRVeGt5bSQC0Wx6vfNLOn9RonDTYjeQySiRd-g_kcj7F3rDH69rIA==", "auth for socket")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/socket.io/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), map[string][]string{"Authorization": []string{*auth}})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		return
		// 	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
