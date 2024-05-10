package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/twitch"
)

func appendItem(msg *twitch.TwitchMessage) []byte {
	return []byte(fmt.Sprintf(`
<turbo-stream action="append" target="messages">
  <template>
    <span>
      %s: %s
    </span>
  </template>
</turbo-stream>
  `, msg.Username, msg.Content))
}

func newListener(streamer string, c *websocket.Conn) twitch.Listener {
	return twitch.Listener{
		Id:       uuid.NewString(),
		Streamer: streamer,
		Callback: func(msg *twitch.TwitchMessage) {
			if err := c.WriteMessage(websocket.TextMessage, appendItem(msg)); err != nil {
				log.Println("write:", err)
			}
		},
	}
}

func WatchStream(app *common.App, c *websocket.Conn) {
	// Add the client to the list of listeners
	listener := newListener(c.Params("id"), c)
	app.Twitch.JoinStreamer(listener)
	c.SetCloseHandler(func(_ int, _ string) error {
		app.Twitch.LeaveStreamer(listener)
		return nil
	})

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
