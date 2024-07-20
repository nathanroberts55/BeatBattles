package cache

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nathanroberts55/beatbattle/soundcloud"
	"github.com/redis/go-redis/v9"
)

type Bucket struct {
	client      *redis.Client
	streamer    string
	cursor      int64
	startCursor int64
}

func NewBucket() *Bucket {
	bucket := &Bucket{
		client:   redisClient(),
		streamer: uuid.NewString(),
	}
	latest := bucket.client.LLen(ctx, bucket.streamer).Val()
	bucket.startCursor = latest
	bucket.cursor = latest
	bucket.client.Expire(ctx, bucket.streamer, time.Hour*4)

	return bucket
}

func (bucket *Bucket) lrange(start, end int64) (result []*soundcloud.EmbededPlayer, err error) {
	res := bucket.client.LRange(ctx, bucket.streamer, start, end)
	if err = res.Err(); err != nil {
		return result, err
	}

	for _, v := range res.Val() {
		var data soundcloud.EmbededPlayer
		err = json.Unmarshal([]byte(v), &data)
		if err != nil {
			return result, err
		}

		result = append(result, &data)
	}

	return result, nil
}

func (bucket *Bucket) PullFromCursor(limit int64) (result []*soundcloud.EmbededPlayer, err error) {
	result, err = bucket.lrange(bucket.cursor, bucket.cursor+limit)
	if err != nil {
		return result, err
	}

	bucket.cursor += int64(len(result))
	return result, err
}

func (bucket *Bucket) Push(sc *soundcloud.EmbededPlayer) error {
	data, err := json.Marshal(sc)
	if err != nil {
		return err
	}

	return bucket.client.RPush(ctx, bucket.streamer, data).Err()
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

func (bucket *Bucket) getThusFar() ([]*soundcloud.EmbededPlayer, error) {
	if bucket.cursor == 0 {
		return []*soundcloud.EmbededPlayer{}, nil
	}

	return bucket.lrange(bucket.startCursor, max(bucket.cursor-1, 0))
}

func (bucket *Bucket) PullUnique(limit int64) (result []*soundcloud.EmbededPlayer, err error) {
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
