package cache

import (
	"encoding/json"

	"github.com/nathanroberts55/beatbattle/soundcloud"
	"github.com/redis/go-redis/v9"
)

type Bucket struct {
	client      *redis.Client
	streamer    string
	cursor      int64
	startCursor int64
}

func NewBucket(streamer string) *Bucket {
	return &Bucket{
		client:   redisClient(),
		streamer: streamer,
	}
}

func (bucket *Bucket) Pull(cursor, limit int64) (result []*soundcloud.SoundcloudItem) {
	data := bucket.client.LRange(ctx, bucket.streamer, cursor, cursor+limit)

	for _, v := range data.Val() {
		var data soundcloud.SoundcloudItem
		err := json.Unmarshal([]byte(v), &data)
		// TODO(dylhack): Handle error
		if err != nil {
			continue
		}

		result = append(result, &data)
	}

	bucket.cursor += int64(len(result))
	return result
}

func (bucket *Bucket) PullFromCursor(limit int64) []*soundcloud.SoundcloudItem {
	if bucket.cursor == 0 {
		latest := bucket.GetLatestCursor()
		bucket.cursor = latest
		bucket.startCursor = latest
	}

	return bucket.Pull(bucket.cursor, limit)
}

func (bucket *Bucket) GetLatestCursor() int64 {
	return bucket.client.LLen(ctx, bucket.streamer).Val()
}

func (bucket *Bucket) Push(sc *soundcloud.SoundcloudItem) error {
	data, err := json.Marshal(sc)
	if err != nil {
		return err
	}

	return bucket.client.LPush(ctx, bucket.streamer, data).Err()
}

func (bucket *Bucket) getThusFar() []*soundcloud.SoundcloudItem {
	return bucket.Pull(bucket.startCursor, bucket.cursor)
}

func (bucket *Bucket) PullUnique(limit int64) (result []*soundcloud.SoundcloudItem) {
	thusFar := bucket.getThusFar()
	for {
		items := bucket.PullFromCursor(limit)

	outter:
		for _, v := range items {
			for _, u := range thusFar {
				if v.Id == u.Id {
					continue outter
				}
			}

			result = append(result, v)
		}

		limit32 := int(limit)
		if len(result) == limit32 || len(items) < limit32 {
			break
		}
	}

	return result
}
