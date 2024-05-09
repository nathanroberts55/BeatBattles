package controllers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/twitch"
)

func newListener(streamer string, c *websocket.Conn) twitch.Listener {
	return twitch.Listener{
		Id:       uuid.NewString(),
		Streamer: streamer,
		Callback: func(msg *twitch.TwitchMessage) {
			data, err := json.Marshal(msg)

			if err != nil {
				log.Println("Failed to serialize FUCK", err)
				return
			}

			if err = c.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
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
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
