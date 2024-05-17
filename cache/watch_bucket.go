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
	bucket := &Bucket{
		client:   redisClient(),
		streamer: streamer,
	}
	latest := bucket.client.LLen(ctx, bucket.streamer).Val()
	bucket.startCursor = latest
	bucket.cursor = latest

	return bucket
}

func (bucket *Bucket) lrange(start, end int64) (result []*soundcloud.SoundcloudItem, err error) {
	res := bucket.client.LRange(ctx, bucket.streamer, start, end)
	if err = res.Err(); err != nil {
		return result, err
	}

	for _, v := range res.Val() {
		var data soundcloud.SoundcloudItem
		err = json.Unmarshal([]byte(v), &data)
		if err != nil {
			return result, err
		}

		result = append(result, &data)
	}

	return result, nil
}

func (bucket *Bucket) PullFromCursor(limit int64) (result []*soundcloud.SoundcloudItem, err error) {
	result, err = bucket.lrange(bucket.cursor, bucket.cursor+limit)
	if err != nil {
		return result, err
	}

	bucket.cursor += int64(len(result))
	return result, err
}

func (bucket *Bucket) Push(sc *soundcloud.SoundcloudItem) error {
	data, err := json.Marshal(sc)
	if err != nil {
		return err
	}

	return bucket.client.LPush(ctx, bucket.streamer, data).Err()
}

func max(a, b int64) int64 {
	if b > a {
		return b
	}
	if a > b {
		return a
	}

	return a
}

func (bucket *Bucket) getThusFar() ([]*soundcloud.SoundcloudItem, error) {
	if bucket.cursor == 0 {
		return []*soundcloud.SoundcloudItem{}, nil
	}

	return bucket.lrange(bucket.startCursor, max(bucket.cursor-1, 0))
}

func (bucket *Bucket) PullUnique(limit int64) (result []*soundcloud.SoundcloudItem, err error) {
	thusFar, err := bucket.getThusFar()
	if err != nil {
		return result, err
	}

	for {
		items, err := bucket.PullFromCursor(limit)
		if err != nil {
			return result, err
		}

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

	return result, nil
}
