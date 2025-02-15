package web

import (
	"encoding/gob"
	"fmt"
	"log/slog"
	"time"

	"wokkibot/handlers"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
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
	version     string
}

func NewServer(config OAuthConfig, bot *wokkibot.Wokkibot, h *handlers.Handler, version string) *Server {
	templatePath := "./web/views"
	if version != "dev" {
		templatePath = "/app/web/views"
	}

	engine := html.New(templatePath, ".html")
	store := session.New(session.Config{
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: true,
	})

	app := fiber.New(fiber.Config{
		Views:        engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	auth := NewAuthHandler(store, config)
	admin := NewAdminHandler(bot, h)

	server := &Server{
		app:         app,
		auth:        auth,
		store:       store,
		oauthConfig: config,
		admin:       admin,
		version:     version,
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

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
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

	GlobalCommands, GuildCommands, err := s.GetCommands()
	if err != nil {
		slog.Error("error getting commands for web view", "error", err)
	}

	return c.Render("dashboard", fiber.Map{
		"User":           userMap,
		"Version":        s.admin.bot.Version,
		"Uptime":         s.GetUptime(),
		"Presence":       s.GetPresence(),
		"Guilds":         s.GetGuildsCount(),
		"GlobalCommands": GlobalCommands,
		"GuildCommands":  GuildCommands,
	})
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) GetPresence() string {
	activity := s.admin.bot.Client.Gateway().Presence().Activities[0]

	switch activity.Type {
	case 0: // Playing
		return fmt.Sprintf("Playing %s", activity.Name)
	case 1: // Streaming
		return fmt.Sprintf("Streaming %s", activity.Name)
	case 2: // Listening
		return fmt.Sprintf("Listening to %s", activity.Name)
	case 3: // Watching
		return fmt.Sprintf("Watching %s", activity.Name)
	case 4: // Custom
		return activity.Name
	case 5: // Competing
		return fmt.Sprintf("Competing in %s", activity.Name)
	default:
		return activity.Name
	}
}

func (s *Server) GetUptime() string {
	return time.Since(s.admin.bot.StartTime).Round(time.Second).String()
}

func (s *Server) GetGuildsCount() int {
	return s.admin.bot.Client.Caches().GuildsLen()
}

func (s *Server) GetCommands() ([]discord.ApplicationCommand, []discord.ApplicationCommand, error) {
	GlobalCommands, err := s.admin.bot.Client.Rest().GetGlobalCommands(s.admin.bot.Client.ApplicationID(), false)
	if err != nil {
		return nil, nil, err
	}

	var GuildCommands []discord.ApplicationCommand
	if s.admin.bot.Config.GuildID != "" {
		GuildCommands, err = s.admin.bot.Client.Rest().GetGuildCommands(s.admin.bot.Client.ApplicationID(), snowflake.MustParse(s.admin.bot.Config.GuildID), false)
		if err != nil {
			return nil, nil, err
		}
	}

	return GlobalCommands, GuildCommands, nil
}
