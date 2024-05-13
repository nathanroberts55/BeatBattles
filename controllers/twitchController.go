package controllers

import (
	"bytes"
	"log"
	"reflect"
	"text/template"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nathanroberts55/beatbattle/cache"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/soundcloud"
	"github.com/nathanroberts55/beatbattle/twitch"
)

type scProps struct {
	IFrame string
}

var pull = []byte("PULL")

func renderEmbed(embed string) []byte {
	props := scProps{
		IFrame: embed,
	}
	var data bytes.Buffer
	tmpl, _ := template.ParseFiles("./views/watch/_embedPlayer.html")
	tmpl.Execute(&data, props)

	return data.Bytes()
}

func newListener(streamer string, bucket *cache.Bucket) twitch.Listener {
	return twitch.Listener{
		Id:       uuid.NewString(),
		Streamer: streamer,
		Callback: func(msg *twitch.TwitchMessage) {
			embed, err := soundcloud.GetEmbed(msg.URL)
			if err != nil {
				log.Printf("Error getting embed for url: '%s' \n | Error:  %v \n", msg.URL, err)
				return
			}

			if err = bucket.Push(renderEmbed(embed)); err != nil {
				log.Printf("Failed to push song to Redis.\n%v\n", err)
			}
		},
	}
}

func WatchStream(app *common.App, c *websocket.Conn) {
	// Add the client to the list of listeners
	streamer := c.Params("id")

	// subscribe
	bucket := cache.NewBucket(streamer)
	listener := newListener(streamer, &bucket)
	app.Twitch.JoinStreamer(listener)

	// respond
	c.SetCloseHandler(func(_ int, _ string) error {
		app.Twitch.LeaveStreamer(listener)
		return nil
	})

	var (
		mt      int
		msg     []byte
		err     error
		payload [][]byte
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if reflect.DeepEqual(msg, pull) {
			payload = bucket.PullFromCursor(20)
		}

		for _, v := range payload {
			if err = c.WriteMessage(mt, v); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
