package twitch

import (
	"errors"
	"log"
	"regexp"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchMessage struct {
	Username string
	Content  string
	Channel  string
	URL      string
}

func extractSong(message string) string {
	re := regexp.MustCompile(`(https:\/\/)(|on\.)(soundcloud\.com)([^ \n]+)`)
	match := re.FindString(message)
	return match
}

func (ts *TwitchService) handleMessage(message twitch.PrivateMessage) (*TwitchMessage, error) {
	url := extractSong(message.Message)
	if url == "" {
		return nil, errors.New("no sound cloud link in message")
	}

	log.Printf("Extacted URL: %s\n", url)
	TwitchMsg := &TwitchMessage{
		Username: message.User.Name,
		Content:  message.Message,
		Channel:  message.Channel,
		URL:      url,
	}

	return TwitchMsg, nil
}
