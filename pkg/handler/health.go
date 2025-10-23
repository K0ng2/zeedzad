package handler

import "github.com/gofiber/fiber/v3"

// DatabaseHealth godoc
// @Summary Database Health check
// @Description Check the Database health status
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} model.DatabaseHealth
// @Router /databasez [get]
func (h *Handler) DatabaseHealth(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	return DatabaseHealth(h.repo.Ping(ctx), c)
}
