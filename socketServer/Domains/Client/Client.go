package Client

import (
	"bytes"
	"log"
	"net/http"
	"socket/socketServer/Domains/Repository/Mongodb"
	model "socket/socketServer/Model"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientInterface struct {
	Client *model.Client
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (clientInterface *ClientInterface) ReadPump() {
	c := clientInterface.Client
	defer func() {
		close(clientInterface.Client.Unregister)
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("%s", message)
		// c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (clientInterface *ClientInterface) WritePump() {
	c := clientInterface.Client
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {

		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		}
	}
}

func ServeWs(hub *model.Hub, db *mgo.Session, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	client := &ClientInterface{
		Client: &model.Client{Conn: c, Send: make(chan []byte, 256), Hub: hub, Unregister: make(chan bool), UserId: bson.ObjectIdHex(r.Header.Get("userId"))},
	}

	hub.RegisterClient <- client.Client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	// find all auctions that I do a bid
	auctions, err := Mongodb.GetAuctionsThatIBid(bson.ObjectIdHex(r.Header.Get("userId")), db)
	if err != nil {
		log.Print("find first auctions: ", err)
		return
	}

	for _, auction := range auctions {
		hub.EnterRoom <- &model.EnterRoom{AuctionId: auction["auctionId"].(bson.ObjectId).Hex(), UserId: client.Client.UserId.Hex()}
	}

	// firstAvgAuctions, err := json.Marshal(avgAuctions)
	// if err != nil {
	// 	log.Print("Error parsing first avg auctions: ", err)
	// 	return
	// }
	// c.WriteMessage(1, firstAvgAuctions)

	<-client.Client.Unregister

}
