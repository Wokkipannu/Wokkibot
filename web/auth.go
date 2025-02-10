package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"log"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AdminUserIDs []snowflake.ID
}

type AuthHandler struct {
	store       *session.Store
	oauthConfig OAuthConfig
}

func NewAuthHandler(store *session.Store, config OAuthConfig) *AuthHandler {
	return &AuthHandler{
		store:       store,
		oauthConfig: config,
	}
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	state := generateRandomState()
	sess, err := h.store.Get(c)
	if err != nil {
		log.Printf("session error in login: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	sess.Set("state", state)
	if err := sess.Save(); err != nil {
		log.Printf("failed to save session in login: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
	}

	log.Printf("login: generated state: %s", state)

	params := url.Values{
		"client_id":     {h.oauthConfig.ClientID},
		"redirect_uri":  {h.oauthConfig.RedirectURI},
		"response_type": {"code"},
		"scope":         {"identify"},
		"state":         {state},
	}

	discordAuthURL := "https://discord.com/api/oauth2/authorize?" + params.Encode()
	return c.Redirect(discordAuthURL)
}

func (h *AuthHandler) HandleCallback(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		log.Printf("session error in callback: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	state := c.Query("state")
	savedState := sess.Get("state")

	log.Printf("callback: comparing states: received: %s, saved: %v", state, savedState)

	if state == "" || state != savedState {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid state")
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("No code provided")
	}

	token, err := h.exchangeCode(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange code: " + err.Error())
	}

	user, err := h.getUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info: " + err.Error())
	}

	isAdmin := false
	for _, adminID := range h.oauthConfig.AdminUserIDs {
		if snowflake.MustParse(user.ID) == adminID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).SendString("Not authorized")
	}

	userData := map[string]string{
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
	}

	sess.Set("user", userData)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session: " + err.Error())
	}

	return c.Redirect("/admin/dashboard")
}

func (h *AuthHandler) exchangeCode(code string) (string, error) {
	data := url.Values{
		"client_id":     {h.oauthConfig.ClientID},
		"client_secret": {h.oauthConfig.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {h.oauthConfig.RedirectURI},
	}

	resp, err := http.PostForm("https://discord.com/api/oauth2/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (h *AuthHandler) getUserInfo(token string) (*DiscordUser, error) {
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *AuthHandler) RequireAuth(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	userData := sess.Get("user")
	if userData == nil {
		return c.Redirect("/login")
	}

	c.Locals("user", userData)
	return c.Next()
}
