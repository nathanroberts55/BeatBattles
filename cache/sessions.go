package cache

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	id     string
	client redis.Client
}

func (s *Session) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.HSet(ctx, s.id, key, data).Err()
}

func (s *Session) Get(key string, out any) error {
	res := s.client.HGet(ctx, s.id, key)
	if res.Err() != nil {
		return res.Err()
	}

	return json.Unmarshal([]byte(res.Val()), out)
}

func GetSession(ctx *fiber.Ctx) *Session {
	id := ctx.Cookies("SESSION_ID", "")
	if len(id) == 0 {
		id = uuid.NewString()
		ctx.Cookie(&fiber.Cookie{
			Name:    "SESSION_ID",
			Value:   id,
			Expires: time.Now().Add(time.Hour),
			Secure:  true,
		})
	}

	return &Session{
		id,
		*redisClient(),
	}
}
