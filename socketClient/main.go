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
var auth = flag.String("Authorization", "Bearer WQzSqFQfDycXmOAU6tk4aCrMVaDEL2WMp6Cc3LQp3-TuH2WWSFhRmSGk0hVRkGaOGLKrsLgmgVmtvcPqjg8Y4-F7vIp1LIdn8DE_xmFQuD4CazBSuPc2Jnk9zE8n71JX8ScuvNH8sF3x9Tm-tYVpNB-RDN8FGIxUawbtxuzqUoOdyEIIrV45P6Ebl49fkLP8zoa_WnhibvRJRT0x9yvpHJWbkkXC823wP62j_Ejk7y6eyjGK1WhHm7zO4rMBIEMSpZZ6-vWkqRUMdFd0KJuaOmJeuw51cDmVWRHk9969cCKGDXNacw==", "auth for socket")

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
