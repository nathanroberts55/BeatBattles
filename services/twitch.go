package services

import (
	"log"
	"strings"
	"sync"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nathanroberts55/beatbattle/helpers"
)

type TwitchService struct {
	client *twitch.Client
	mu     sync.Mutex
}
type TwitchMessage struct {
	Username string
	Content  string
	Channel  string
	URL      string
}

func NewTwitchService() *TwitchService {
	return &TwitchService{
		client: twitch.NewAnonymousClient(),
	}
}

func HandleMessage(message twitch.PrivateMessage) (*TwitchMessage, error) {
	if !strings.Contains(message.Message, "soundclo") {
		return nil, nil
	}

	url := helpers.ExtractSong(message.Message)
	if url != "" {
		log.Printf("SoundCloud URL: %s \n", url)
	}

	return &TwitchMessage{
		Username: message.User.Name,
		Content:  message.Message,
		Channel:  message.Channel,
		URL:      url,
	}, nil
}

func (s *TwitchService) ConnectToTwitch() {
	s.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msg, err := HandleMessage(message)
		if err != nil {
			log.Printf("Unable to Handle Message: %s", err)
		}

		log.Println(msg)
	})

	err := s.client.Connect()
	if err != nil {
		panic(err)
	}
}

func (s *TwitchService) JoinStreamer(streamerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.client.Join(streamerName)
}
