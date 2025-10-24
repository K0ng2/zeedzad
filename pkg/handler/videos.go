package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/K0ng2/zeedzad/model"
)

// GetVideos godoc
// @Summary Get all videos
// @Description Get all videos with optional search and pagination
// @Tags videos
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(24)
// @Param search query string false "Search by video title or game name"
// @Success 200 {object} model.APIResponse[[]model.VideoResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /videos [get]
func (h *Handler) GetVideos(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	q, err := GetOffset(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(model.Error{Error: err.Error()})
	}

	search := c.Query("search", "")

	videos, err := h.repo.GetVideos(ctx, *q, search)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	total, err := h.repo.GetVideoTotalItems(ctx, search)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	meta := &model.Meta{
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}

	return c.JSON(Response(videos, meta))
}

// GetVideoByID godoc
// @Summary Get video by ID
// @Description Get a single video by its ID
// @Tags videos
// @Accept  json
// @Produce  json
// @Param id path string true "Video ID"
// @Success 200 {object} model.APIResponse[model.VideoResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /videos/{id} [get]
func (h *Handler) GetVideoByID(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(ErrInvalidPathParams)
	}

	video, err := h.repo.GetVideoByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	return c.JSON(Response(video, nil))
}

// UpdateVideoGame godoc
// @Summary Update video's game match
// @Description Match a video with a game
// @Tags videos
// @Accept  json
// @Produce  json
// @Param id path string true "Video ID"
// @Param game body model.UpdateVideoGameRequest true "Game ID to match"
// @Success 200
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /videos/{id}/game [put]
func (h *Handler) UpdateVideoGame(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	videoID := c.Params("id")
	if videoID == "" {
		return c.Status(http.StatusBadRequest).JSON(ErrInvalidPathParams)
	}

	var requestBody model.UpdateVideoGameRequest
	if err := c.Bind().JSON(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrInvalidRequestBody)
	}

	err := h.repo.UpdateVideoGame(ctx, videoID, requestBody.GameID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}
