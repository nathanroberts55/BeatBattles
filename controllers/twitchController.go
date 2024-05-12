package controllers

import (
	"bytes"
	"log"
	"text/template"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/soundcloud"
	"github.com/nathanroberts55/beatbattle/twitch"
)

type scProps struct {
	IFrame string
}

func appendItem(embed string) []byte {
	props := scProps{
		IFrame: embed,
	}
	var data bytes.Buffer

	tmpl, _ := template.ParseFiles("./views/partials/embedPlayer.html")

	tmpl.Execute(&data, props)

	return data.Bytes()
}

func newListener(streamer string, c *websocket.Conn) twitch.Listener {
	return twitch.Listener{
		Id:       uuid.NewString(),
		Streamer: streamer,
		Callback: func(msg *twitch.TwitchMessage) {
			embed, err := soundcloud.GetEmbed(msg.URL)
			if err != nil {
				log.Printf("Error getting embed for url: '%s' \n | Error:  %v \n", msg.URL, err)
				return
			}
			if err := c.WriteMessage(websocket.TextMessage, appendItem(embed)); err != nil {
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
