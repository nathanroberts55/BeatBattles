package watch

import (
	"log"

	"github.com/nathanroberts55/beatbattle/cache"
	"github.com/nathanroberts55/beatbattle/common"
	"github.com/nathanroberts55/beatbattle/controllers"
	"regexp"
)

func remember(ctx *common.Ctx, streamer string) error {
	recent, err := recentStreamers(ctx)
	if err != nil {
		return err
	}

	for _, v := range recent {
		if v == streamer {
			return nil
		}
	}

	if len(recent) == 10 {
		recent = append(recent[:9], streamer)
	} else {
		recent = append(recent, streamer)
	}

	store := cache.GetSession(ctx.Ctx)
	return store.Set(STREAMER_STORE, streamerStore{recent})
}

// /watch/:streamer
type showProps struct {
	controllers.Params
	Streamer string
}

func sanatize(streamer string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")
	clean := regex.ReplaceAllString(streamer, "")
	if len(clean) > 25 {
		clean = clean[:25]
	}

	return clean
}

func Show(c *common.Ctx) error {
	streamer := sanatize(c.Params("streamer", "ttlnow"))
	if err := remember(c, streamer); err != nil {
		log.Printf("Failed to remember streamer\n%s\n", err)
	}

	return c.Render("watch/show", showProps{
		controllers.DefaultParams,
		streamer,
	}, "layouts/watch")
}
