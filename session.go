package session

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/storage/redis"
)

var rememberMeExpiry = 7 * 24 * time.Hour
var defaultExpiry = 30 * time.Minute

var DefaultSession = session.New(session.Config{
	Expiration:     defaultExpiry,
	CookieName:     "Verify-Session",
	CookieHTTPOnly: true,
	Storage:        memory.New(),
})

var RememberMeSession = session.New(session.Config{
	Expiration:     rememberMeExpiry,
	CookieName:     "Verify-Session-Remember",
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

func Default(cfg Config) {
	DefaultSession = New(cfg)
}
func RememberMe(cfg Config) {
	RememberMeSession = New(cfg)
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
	case "memory":
		store = memory.New()
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
	return session.New(session.Config{
		Expiration:     cfg.Expiration,
		CookieName:     cfg.CookieName,
		CookieHTTPOnly: cfg.CookieHttpOnly,
		Storage:        store,
	})
}

func SetKeys(c *fiber.Ctx, data fiber.Map) error {
	sess := mustPickSession(c)
	for key, value := range data {
		sess.Set(key, value)
	}
	return sess.Save()
}

func Delete(c *fiber.Ctx, key string) error {
	sess := mustPickSession(c)
	sess.Delete(key)
	return sess.Save()
}

func DeleteKeys(c *fiber.Ctx, keys ...string) error {
	sess := mustPickSession(c)
	for _, key := range keys {
		sess.Delete(key)
	}
	return sess.Save()
}

func DeleteWithDestroy(c *fiber.Ctx, keys ...string) error {
	sess := mustPickSession(c)
	for _, key := range keys {
		sess.Delete(key)
	}
	Destroy(c)
	return sess.Save()
}

func Get(c *fiber.Ctx, key string) (interface{}, error) {
	sess := mustPickSession(c)
	return sess.Get(key), nil
}

func Destroy(c *fiber.Ctx) error {
	sess := mustPickSession(c)
	sess.Destroy()
	return sess.Save()
}

func Save(c *fiber.Ctx) error {
	sess := mustPickSession(c)
	return sess.Save()
}

func Fresh(c *fiber.Ctx) (bool, error) {
	sess := mustPickSession(c)
	return sess.Fresh(), nil
}

func ID(c *fiber.Ctx) (string, error) {
	sess := mustPickSession(c)
	return sess.ID(), nil
}

func Regenerate(c *fiber.Ctx) error {
	sess := mustPickSession(c)
	return sess.Regenerate()
}

func mustPickSession(c *fiber.Ctx) *session.Session {
	sess, err := DefaultSession.Get(c)
	if err != nil {
		panic(err)
	}
	return sess
}
