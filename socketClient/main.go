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

var addr = flag.String("addr", "localhost:5000", "http service address")
var auth = flag.String("Authorization", "Bearer _QM8B8Yqj_UHOVWGKW-vJ2E4Dn_BwOWuuAdowFXovZq3oGEaFNwkm0Ns9Az8ULDcb34Z_FxbYW9xtfeQylHr_YaDpYeodAOPTFewiLI3sPVTxfYucreBI14Sd_t92HuyPZlsaJ9V9eDcBWf_wXuB7yTd7flJ2B7f3hlWoOz3RuMTxe5fvlIHwYBNDXGAxb1mOjU9g9ieMYsmAQ1Y-SZemS6GXw_uLvk8aG46ZY2tXWirGH4Fobcb6kQAUUdLZqfyEuAkjqpHNhHCyXLry0U0vCWMuSBgm2AedpHOT1dMOLYV4Cg=", "auth for socket")

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
