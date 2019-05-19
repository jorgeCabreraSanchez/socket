package main

import (
	"flag"
	"log"
	"time"

	"github.com/gorilla/websocket"
	socketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	auth := flag.String("Authorization", "Bearer _qM8B8Yqj_UHOVWGKW-vJ2E4Dn_BwOWuuAdowFXovZq3oGEaFNwkm0Ns9Az8ULDcb34Z_FxbYW9xtfeQylHr_YaDpYeodAOPTFewiLI3sPVTxfYucreBI14Sd_t92HuyPZlsaJ9V9eDcBWf_wXuB7yTd7flJ2B7f3hlWoOz3RuMTxe5fvlIHwYBNDXGAxb1mOjU9g9ieMYsmAQ1Y-SZemS6GXw_uLvk8aG46ZY2tXWirGH4Fobcb6kQAUUdLZqfyEuAkjqpHNhHCyXLry0U0vCWMuSBgm2AedpHOT1dMOLYV4Cg=", "auth for socket")
	flag.Parse()

	url := socketio.GetUrl("localhost", 5000, false)
	defaultTransport := transport.GetDefaultWebsocketTransport()

	defaultTransport.RequestHeader = make(map[string][]string)
	defaultTransport.RequestHeader.Set("Authorization", *auth)

	log.Print(*auth)

	socket, err := socketio.Dial(url, defaultTransport)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			log.Print(err)
		} else {
			log.Fatal(err)
		}
	} else {
		defer socket.Close()

		socket.On("message", func(c *socketio.Channel, tal string) {
			log.Print(tal)
		})

		type matchBid struct {
			AuctionId bson.ObjectId `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
			UserId    bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
		}

		socket.On("connection", func(c *socketio.Channel) {
			err := c.Emit("subscribteToAnAuction", "hola")
			if err != nil {
				log.Fatal(err)
			}

		})

		<-time.After(time.Duration(24 * time.Hour))
	}

}
