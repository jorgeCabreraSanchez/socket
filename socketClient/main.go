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
var auth = flag.String("Authorization", "Bearer s0aKwZJ2JRmO9n8O-ZIJN0iIH-zpHfAeJ9WHmgAMtaD2s_JsqquzX2l2oAE-kOagXkNqOrbwNps70NEY1PKvNTilRA_aKxtM2VwUdw8Rb5fORvIzbrP07o0p9YsvLpbimh-7ugMAcZmmqytyQyeW5R9JLmW2LomR_2jxJNiDInARNsGKa34Mw3s5ASH7Kjodqwu_63GOzpHRjRJLnZHgOlNV6SDqzLP0rnkYT2y_5520Te9-NQiZct_pKPlIiKx2Z0oy78c9Wj8_zrUjG2ZtOWltrhomZ526DtHCOsC9X70rWp0YXg==", "auth for socket")

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
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
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
