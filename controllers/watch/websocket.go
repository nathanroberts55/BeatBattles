package watch

import (
	"log"
	"reflect"

	"bytes"
	"text/template"

	"github.com/gofiber/contrib/websocket"
	"github.com/nathanroberts55/beatbattle/cache"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/soundcloud"
	"github.com/nathanroberts55/beatbattle/twitch"
)

type scProps struct {
	IFrame string
}

func renderEmbed(embed string) (out []byte) {
	tmpl, err := template.ParseFiles("./views/watch/_embedPlayer.html")
	if err != nil {
		log.Printf("Failed to parse template.\n%v\n", err)
		return out
	}

	var data bytes.Buffer
	err = tmpl.Execute(&data, scProps{
		IFrame: embed,
	})

	if err != nil {
		log.Printf("Failed to execute template.\n%v\n", err)
		return out
	}

	return data.Bytes()
}

func newListener(streamer string, bucket *cache.Bucket) *twitch.Listener {
	return twitch.NewListener(
		streamer,
		func(msg *twitch.TwitchMessage) {
			sc, err := soundcloud.GetEmbed(msg.URL)
			if err != nil {
				log.Printf("Failed to get oEmbed for url: '%s'\n%v\n", msg.URL, err)
				return
			}

			if err = bucket.Push(sc); err != nil {
				log.Printf("Failed to push song to Redis.\n%v\n", err)
			}
		},
	)
}

func subscribe(app *common.App, c *websocket.Conn) *cache.Bucket {
	streamer := c.Params("streamer")
	bucket := cache.NewBucket(streamer)
	listener := newListener(streamer, bucket)
	app.Twitch.JoinStreamer(listener)

	c.SetCloseHandler(func(_ int, _ string) error {
		app.Twitch.LeaveStreamer(listener)
		return nil
	})

	return bucket
}

// /ws/watch/:streamer
var pullMsg = []byte("PULL")

func Watch(app *common.App, c *websocket.Conn) {
	var (
		err    error
		mt     int
		msg    []byte
		resp   [][]byte
		bucket = subscribe(app, c)
	)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if reflect.DeepEqual(msg, pullMsg) {
			for _, v := range bucket.PullUnique(20) {
				resp = append(resp, renderEmbed(v.Html))
			}
		}

		for _, v := range resp {
			if err = c.WriteMessage(mt, v); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
