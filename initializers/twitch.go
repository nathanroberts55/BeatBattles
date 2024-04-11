package initializers

import (
	"log"

	"github.com/gempir/go-twitch-irc/v4"
)

func ConnectToTwitch(streamerName string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf("%s: %s \n", message.User.DisplayName, message.Message)
	})

	client.Join(streamerName)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
