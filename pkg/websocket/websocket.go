package websocket

import (
	"bytes"
	"context"
	"fmt"
	"github.com/maxchagin/go-memorycache-example"
	"lets-go-chat-v2/internal/auth"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/pkg/logging"
	"lets-go-chat-v2/pkg/utils/cache"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS struct {
	logger *logging.Logger
	cache  *memorycache.Cache
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	UserName string `json:"UserName"`
	// The websocket connection.
	conn   *websocket.Conn
	cache  *memorycache.Cache
	logger *logging.Logger
	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.CreateNewMessage(string(message))
		c.hub.broadcast <- message

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.removeFromCache(c.UserName)
		c.conn.Close()

	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var Cache *memorycache.Cache

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	userCtxValue := r.Context().Value(middleware.UserContextKey)
	if userCtxValue == nil {
		log.Println("Not authenticated")
		return
	}

	user := userCtxValue.(**auth.UserClaims)
	//Cache.Set("webSocketUsers", connCache.(string)+(*user).UserName+"!"+token[0]+":", 10*time.Minute)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, UserName: (*user).UserName, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.sendMissedMessages()
	go client.writePump()
	go client.readPump()
}

func (c *Client) removeFromCache(token string) bool {
	cacheConnections, err := c.cache.Get("webSocketUsers")

	if err && cacheConnections != nil {
		arr := strings.Split(cacheConnections.(string), ":")
		var newCacheConnections []string
		for key, value := range arr {
			if value != "" {
				arrCred := strings.Split(value, "!")
				if token == arrCred[1] {
					newCacheConnections = append(arr[:key], arr[key+1:]...)
					c.cache.Set("webSocketUsers", strings.Join(newCacheConnections, ":"), 10*time.Minute)
					return true
				}
			}
		}
	}
	return false
}

func (c *Client) saveToCache(token string) {
	cacheConnections, err := cache.Cache.Get("webSocketUsers")
	fmt.Println("saveToCache", cacheConnections)
	if !err || cacheConnections == nil {
		cacheConnections = ""
	}
	c.cache.Set("activeUsers", cacheConnections.(string)+c.UserName+"!"+token+":", 10*time.Minute)
}

func (c *Client) sendMissedMessages() {

	userMessages, _ := c.hub.messageRepository.GetUnreadMessages(context.TODO())
	for _, message := range userMessages {
		fmt.Println(message)
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write([]byte(message.Message))
	}

}

func (c *Client) CreateNewMessage(message string) {
	user, _, err := c.hub.userRepository.GetUser(context.TODO(), c.UserName)
	if err != nil {
		c.logger.Error(err)
	}
	_, err = c.hub.messageRepository.CreateUserMessage(context.TODO(), user[0].ID, message)
	if err != nil {
		c.logger.Error(err)
	}

}
