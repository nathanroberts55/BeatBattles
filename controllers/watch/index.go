package watch

import (
	"fmt"
	"log"

	"github.com/nathanroberts55/beatbattle/cache"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/controllers"
)

var STREAMER_STORE = "STREAMERS"

type streamerStore struct {
	Streamers []string `json:"streamers"`
}

func recentStreamers(ctx *common.Ctx) ([]string, error) {
	store := cache.GetSession(ctx.Ctx)
	var recent streamerStore
	store.Get(STREAMER_STORE, &recent)

	return recent.Streamers, nil
}

// /watch
type indexProps struct {
	controllers.Params
	RecentStreamers []string
}

func Index(c *common.Ctx) error {
	streamer := c.Query("streamer", "")
	if len(streamer) > 0 {
		c.Redirect(fmt.Sprintf("/watch/%s", streamer))
		return nil
	}

	recent, err := recentStreamers(c)
	if err != nil {
		log.Printf("Failed to get recent streamers\n%v\n", err)
	}

	return c.Render("watch/index", indexProps{
		controllers.DefaultParams,
		recent,
	})
}
