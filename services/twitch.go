package services

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
}

func HandleMessage(twitchMessage TwitchMessage) error {
	if !strings.Contains(twitchMessage.Content, "soundclo") {
		return nil
	}
	return nil
}

func ExtractSong(message string) string {
	re := regexp.MustCompile(`https://soundcloud\.com/[^/]+/[^/?]+`)
	match := re.FindString(message)
	return match
}

func ConnectToTwitch(streamerName string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		url := ExtractSong(message.Message)
		if url != "" {
			log.Printf("SoundCloud URL: %s \n", url)
		}
	})
	client.Join(streamerName)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
