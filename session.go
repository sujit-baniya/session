package session

import (
	"github.com/gofiber/fiber/v2"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var DefaultSession = session.New(session.Config{
	Expiration:     30 * time.Minute,
	CookieName:     "Verify-Session",
	CookieHTTPOnly: true,
	Storage:        redis.New(),
})

type Config struct {
	Host           string
	Port           int
	DB             int
	Username       string
	Password       string
	Expiration     time.Duration
	CookieName     string
	CookieHttpOnly bool
}

func New(cfg Config) *session.Store {
	DefaultSession = session.New(session.Config{
		Expiration:     cfg.Expiration,
		CookieName:     cfg.CookieName,
		CookieHTTPOnly: cfg.CookieHttpOnly,
		Storage: redis.New(redis.Config{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.DB,
		}),
	})
	return DefaultSession
}

func Set(c *fiber.Ctx, key string, value interface{}) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	store.Set(key, value)
	return nil
}

func Delete(c *fiber.Ctx, key string) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	store.Delete(key)
	return nil
}

func Get(c *fiber.Ctx, key string) (interface{}, error) {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return nil, err
	}
	return store.Get(key), nil
}

func Destroy(c *fiber.Ctx) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	return store.Destroy()
}

func Save(c *fiber.Ctx) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	return store.Save()
}

func Fresh(c *fiber.Ctx) (bool, error) {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return false, err
	}
	return store.Fresh(), nil
}

func ID(c *fiber.Ctx) (string, error) {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return "", err
	}
	return store.ID(), nil
}

func Regenerate(c *fiber.Ctx) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	return store.Regenerate()
}
