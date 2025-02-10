package web

import (
	"wokkibot/database"
	"wokkibot/handlers"
	"wokkibot/types"
	"wokkibot/wokkibot"

	"github.com/disgoorg/snowflake/v2"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	bot      *wokkibot.Wokkibot
	handlers *handlers.Handler
}

func NewAdminHandler(bot *wokkibot.Wokkibot, h *handlers.Handler) *AdminHandler {
	return &AdminHandler{
		bot:      bot,
		handlers: h,
	}
}

func (h *AdminHandler) GetCustomCommands(c *fiber.Ctx) error {
	commands, err := database.GetCustomCommands()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(commands)
}

func (h *AdminHandler) AddCustomCommand(c *fiber.Ctx) error {
	var dbCmd database.CustomCommand
	if err := c.BodyParser(&dbCmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cmd := types.Command{
		GuildID:     snowflake.MustParse(dbCmd.GuildID),
		Prefix:      dbCmd.Prefix,
		Name:        dbCmd.Name,
		Description: dbCmd.Description,
		Output:      dbCmd.Output,
		Author:      snowflake.MustParse(dbCmd.Author),
	}

	if err := h.handlers.AddOrUpdateCommand(cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *AdminHandler) UpdateCustomCommand(c *fiber.Ctx) error {
	var dbCmd database.CustomCommand
	if err := c.BodyParser(&dbCmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cmd := types.Command{
		GuildID:     snowflake.MustParse(dbCmd.GuildID),
		Prefix:      dbCmd.Prefix,
		Name:        dbCmd.Name,
		Description: dbCmd.Description,
		Output:      dbCmd.Output,
		Author:      snowflake.MustParse(dbCmd.Author),
	}

	if err := h.handlers.AddOrUpdateCommand(cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(dbCmd)
}

func (h *AdminHandler) DeleteCustomCommand(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	cmd, err := database.GetCommandByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.handlers.RemoveCommand(cmd.Prefix, cmd.Name, snowflake.MustParse(cmd.Author)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) GetFridayClips(c *fiber.Ctx) error {
	clips, err := database.GetFridayClips()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(clips)
}

func (h *AdminHandler) AddFridayClip(c *fiber.Ctx) error {
	var clip struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&clip); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.AddFridayClip(clip.URL); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *AdminHandler) DeleteFridayClip(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := database.DeleteFridayClip(int64(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AdminHandler) GetPizzaToppings(c *fiber.Ctx) error {
	toppings, err := database.GetPizzaToppings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(toppings)
}

func (h *AdminHandler) AddPizzaTopping(c *fiber.Ctx) error {
	var topping struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&topping); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.AddPizzaTopping(topping.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *AdminHandler) DeletePizzaTopping(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := database.DeletePizzaTopping(int64(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
