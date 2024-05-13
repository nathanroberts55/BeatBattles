package cache

import (
	"log"

	"github.com/redis/go-redis/v9"
)

type Bucket struct {
	client   *redis.Client
	streamer string
	cursor   int64
}

func NewBucket(streamer string) Bucket {
	return Bucket{
		client:   redisClient(),
		streamer: streamer,
	}
}

func (bucket *Bucket) Pull(cursor, limit int64) [][]byte {
	log.Printf("Bucket state {%v}\n", bucket)
	data := bucket.client.LRange(ctx, bucket.streamer, cursor, cursor+limit)
	var result [][]byte

	for _, v := range data.Val() {
		result = append(result, []byte(v))
	}

	bucket.cursor += int64(len(result))
	return result
}

func (bucket *Bucket) PullFromCursor(limit int64) [][]byte {
	if bucket.cursor == 0 {
		return bucket.PullLast(limit)
	}

	return bucket.Pull(bucket.cursor, limit)
}

func (bucket *Bucket) PullLast(limit int64) [][]byte {
	cursor := bucket.client.LLen(ctx, bucket.streamer).Val() - limit
	if cursor < 0 {
		cursor = 0
	}

	return bucket.Pull(cursor, limit)
}

func (bucket *Bucket) Push(embed []byte) error {
	return bucket.client.LPush(ctx, bucket.streamer, embed).Err()
}
