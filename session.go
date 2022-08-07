package session

import (
	"encoding/gob"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/mysql"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/storage/redis"
)

var RememberMeExpiry = 30 * 24 * time.Hour
var DefaultSessionExpiry = 30 * time.Minute

var DefaultSession = session.New(session.Config{
	Expiration:     DefaultSessionExpiry,
	KeyLookup:      "cookie:Verify-Session",
	CookieHTTPOnly: true,
	Storage:        memory.New(),
})

type Config struct {
	Driver         string `yaml:"driver" env:"DB_DRIVER"`
	Host           string `yaml:"host" env:"DB_HOST"`
	Username       string `yaml:"username" env:"DB_USER"`
	Password       string `yaml:"password" env:"DB_PASS"`
	DB             string `yaml:"db" env:"DB_NAME"`
	Port           int    `yaml:"port" env:"DB_PORT"`
	Table          string
	Expiration     time.Duration
	CookieName     string
	CookieHttpOnly bool
	RegisterTypes  []interface{}
}

func Default(cfg Config) {
	DefaultSession = New(cfg)
	for _, i := range cfg.RegisterTypes {
		Register(i)
	}
}

func New(cfg Config) *session.Store {
	var store fiber.Storage
	cfg.CookieHttpOnly = true
	if cfg.Expiration == 0 {
		cfg.Expiration = DefaultSessionExpiry
	}
	if cfg.CookieName == "" {
		cfg.CookieName = "cookie:Verify-Session"
	}
	if cfg.Table == "" {
		cfg.Table = "login_sessions"
	}
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
		break
	case "mysql":
		store = mysql.New(mysql.Config{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.DB,
			Table:    cfg.Table,
		})
		break
	case "memory":
		store = memory.New()
		break
	default:
		db, _ := strconv.Atoi(cfg.DB)
		store = redis.New(redis.Config{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: db,
		})
		break
	}
	return session.New(session.Config{
		Expiration:     cfg.Expiration,
		KeyLookup:      cfg.CookieName,
		CookieHTTPOnly: cfg.CookieHttpOnly,
		Storage:        store,
	})
}

func Set(c *fiber.Ctx, key string, value interface{}, exp ...time.Duration) error {
	sess := mustPickSession(c)
	sess.Set(key, value)
	if len(exp) > 0 {
		sess.SetExpiry(exp[0])
	}
	return sess.Save()
}

func SetKeys(c *fiber.Ctx, data fiber.Map, exp ...time.Duration) error {
	sess := mustPickSession(c)
	for key, value := range data {
		sess.Set(key, value)
	}
	if len(exp) > 0 {
		sess.SetExpiry(exp[0])
	}
	return sess.Save()
}

func Delete(c *fiber.Ctx, key string) error {
	sess := mustPickSession(c)
	sess.Delete(key)
	return sess.Save()
}

func RememberMe(c *fiber.Ctx) error {
	return SetExpiry(c, RememberMeExpiry)
}

func SetExpiry(c *fiber.Ctx, exp time.Duration) error {
	sess := mustPickSession(c)
	sess.SetExpiry(exp)
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

func SetUser(c *fiber.Ctx, user interface{}) error {
	return Set(c, "user", user)
}

func User(c *fiber.Ctx) (interface{}, error) {
	return Get(c, "user")
}

func Register(i interface{}) {
	gob.Register(i)
}
