package twitch

import (
	"log"
	"sync"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchService struct {
	client       *twitch.Client
	Streams      map[string][]*Listener
	StreamsMutex sync.Mutex
}

func New() *TwitchService {
	return &TwitchService{
		client:       twitch.NewAnonymousClient(),
		Streams:      make(map[string][]*Listener),
		StreamsMutex: sync.Mutex{},
	}
}

func (ts *TwitchService) Start() {
	ts.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msg, err := ts.handleMessage(message)
		if err != nil {
			log.Printf("Unable to Handle Message: %s", err)
			return
		}

		if listeners, isOK := ts.Streams[message.Channel]; isOK {
			for _, listener := range listeners {
				listener.Callback(msg)
			}
		}
	})

	err := ts.client.Connect()
	if err != nil {
		panic(err)
	}
}
