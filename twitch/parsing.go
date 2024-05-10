package twitch

import (
	"log"
	"regexp"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchMessage struct {
	Username string
	Content  string
	Channel  string
	URL      string
}

func extractSong(message string) string {
	re := regexp.MustCompile(`https://soundcloud\.com/[^/]+/[^/?]+`)
	match := re.FindString(message)
	return match
}

func (ts *TwitchService) handleMessage(message twitch.PrivateMessage) (*TwitchMessage, error) {
	if !strings.Contains(message.Message, "soundclo") {
		return nil, nil
	}

	url := extractSong(message.Message)
	if url != "" {
		log.Printf("SoundCloud URL: %s \n", url)
	}

	TwitchMsg := &TwitchMessage{
		Username: message.User.Name,
		Content:  message.Message,
		Channel:  message.Channel,
		URL:      url,
	}

	return TwitchMsg, nil
}

func (ts *TwitchService) sampleHandleMessage(message twitch.PrivateMessage) (*TwitchMessage, error) {
	TwitchMsg := &TwitchMessage{
		Username: message.User.Name,
		Content:  message.Message,
		Channel:  message.Channel,
	}

	return TwitchMsg, nil
}
