package server

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

func NewSessionStore() *session.Store {
	// TODO(dylhack): configure
	return session.New(session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:session_id",
		KeyGenerator: utils.UUIDv4,
	})
}
