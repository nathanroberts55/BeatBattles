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
			} else {
				log.Printf("Pushed song to Redis: %s\n", sc.Id)
			}
		},
	)
}

func subscribe(app *common.App, c *websocket.Conn) *cache.Bucket {
	streamer := c.Params("streamer")
	bucket := cache.NewBucket()
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
	streamer := c.Params("streamer")
	log.Printf("New connection: %s\n", streamer)

	var (
		err    error
		mt     int
		msg    []byte
		resp   [][]byte
		bucket = subscribe(app, c)
	)

	reset := func() {
		resp = [][]byte{}
	}

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if reflect.DeepEqual(msg, pullMsg) {
			items, err := bucket.PullUnique(20)
			if err != nil {
				log.Printf("Failed to pull items from Redis.\n%v\n", err)
				continue
			}

			for _, v := range items {
				log.Printf("%s | %s: pulled %s\n", c.Cookies("SESSION_ID"), streamer, v.Id)
				resp = append(resp, renderEmbed(v.Html))
			}
		}

		if len(resp) == 0 {
			log.Println("No items pulled.")
		}

		for _, v := range resp {
			if err = c.WriteMessage(mt, v); err != nil {
				log.Println("write:", err)
				break
			}
		}

		reset()
	}
}
