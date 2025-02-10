package web

import (
	"encoding/gob"
	"fmt"
	"time"

	"wokkibot/handlers"
	"wokkibot/wokkibot"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

func init() {
	gob.Register(map[string]string{})
}

type Server struct {
	app         *fiber.App
	auth        *AuthHandler
	store       *session.Store
	oauthConfig OAuthConfig
	admin       *AdminHandler
}

func NewServer(config OAuthConfig, bot *wokkibot.Wokkibot, h *handlers.Handler) *Server {
	engine := html.New("/app/web/views", ".html")
	store := session.New(session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:session_id",
		CookiePath:   "/",
		CookieSecure: true,
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	auth := NewAuthHandler(store, config)
	admin := NewAdminHandler(bot, h)

	server := &Server{
		app:         app,
		auth:        auth,
		store:       store,
		oauthConfig: config,
		admin:       admin,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.app.Use(logger.New())

	s.app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("Request: %s %s\n", c.Method(), c.Path())
		return c.Next()
	})

	s.app.Get("/login", s.auth.HandleLogin)
	s.app.Get("/callback", s.auth.HandleCallback)

	admin := s.app.Group("/admin", s.auth.RequireAuth)

	admin.Get("/dashboard", s.handleDashboard)

	api := admin.Group("/api")

	api.Get("/commands", s.admin.GetCustomCommands)
	api.Post("/commands", s.admin.AddCustomCommand)
	api.Put("/commands/:id", s.admin.UpdateCustomCommand)
	api.Delete("/commands/:id", s.admin.DeleteCustomCommand)

	api.Get("/friday-clips", s.admin.GetFridayClips)
	api.Post("/friday-clips", s.admin.AddFridayClip)
	api.Delete("/friday-clips/:id", s.admin.DeleteFridayClip)

	api.Get("/pizza-toppings", s.admin.GetPizzaToppings)
	api.Post("/pizza-toppings", s.admin.AddPizzaTopping)
	api.Delete("/pizza-toppings/:id", s.admin.DeletePizzaTopping)
}

func (s *Server) handleDashboard(c *fiber.Ctx) error {
	sess, err := s.store.Get(c)
	if err != nil {
		return c.Redirect("/login")
	}

	userData := sess.Get("user")
	if userData == nil {
		return c.Redirect("/login")
	}

	userMap, ok := userData.(map[string]string)
	if !ok {
		sess.Delete("user")
		sess.Save()
		return c.Redirect("/login")
	}

	avatarURL := ""
	if userMap["avatar"] != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userMap["id"], userMap["avatar"])
	}
	userMap["avatar_url"] = avatarURL

	return c.Render("dashboard", fiber.Map{
		"User": userMap,
	})
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}
