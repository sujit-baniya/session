package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/postgres"
	"strconv"
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
	Driver         string
	Host           string
	Port           int
	DB             string
	Table          string
	Username       string
	Password       string
	Expiration     time.Duration
	CookieName     string
	CookieHttpOnly bool
}

func New(cfg Config) *session.Store {
	var store fiber.Storage
	switch cfg.Driver {
	case "postgres":
		store = postgres.New(postgres.Config{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.DB,
			Table:    cfg.Table,
		})
	default:
		db, _ := strconv.Atoi(cfg.DB)
		store = redis.New(redis.Config{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: db,
		})
	}
	DefaultSession = session.New(session.Config{
		Expiration:     cfg.Expiration,
		CookieName:     cfg.CookieName,
		CookieHTTPOnly: cfg.CookieHttpOnly,
		Storage:        store,
	})
	return DefaultSession
}

func SetKeys(c *fiber.Ctx, data fiber.Map) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	for key, value := range data {
		store.Set(key, value)
	}
	return store.Save()
}

func Delete(c *fiber.Ctx, key string) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	store.Delete(key)
	return store.Save()
}

func DeleteKeys(c *fiber.Ctx, keys ...string) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	for _, key := range keys {
		store.Delete(key)
	}
	return store.Save()
}

func DeleteWithDestroy(c *fiber.Ctx, keys ...string) error {
	store, err := DefaultSession.Get(c)
	if err != nil {
		return err
	}
	for _, key := range keys {
		store.Delete(key)
	}
	Destroy(c)
	return store.Save()
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
	store.Destroy()
	return store.Save()
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
