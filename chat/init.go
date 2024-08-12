package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
)

func Initialize(e *gin.Engine) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "chrome-extension://oilioclnckkoijghdniegedkbocfpnip"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	e.Use(c)

	hub := NewHub()

	e.GET("/chat/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		go hub.ConnectClient(conn)
	})
}
