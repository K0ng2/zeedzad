package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v3"

	"github.com/K0ng2/zeedzad/model"
	"github.com/K0ng2/zeedzad/utils"
)

// GetGames godoc
// @Summary Get all games
// @Description Get all games with optional search and pagination
// @Tags games
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(20)
// @Param search query string false "Search by game name"
// @Success 200 {object} model.APIResponse[[]model.GameResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /games [get]
func (h *Handler) GetGames(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	q, err := GetOffset(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(model.Error{Error: err.Error()})
	}

	search := c.Query("search", "")

	games, err := h.repo.GetGames(ctx, *q, search)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	total, err := h.repo.GetGameTotalItems(ctx, search)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	meta := &model.Meta{
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}

	return c.JSON(Response(games, meta))
}

// GetGameByID godoc
// @Summary Get game by ID
// @Description Get a single game by its ID
// @Tags games
// @Accept  json
// @Produce  json
// @Param id path int true "Game ID"
// @Success 200 {object} model.APIResponse[model.GameResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /games/{id} [get]
func (h *Handler) GetGameByID(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	id, err := fiber.Convert(c.Params("id"), utils.Atoi64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrInvalidPathParams)
	}

	game, err := h.repo.GetGameByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	return c.JSON(Response(game, nil))
}

// SearchSteamGames godoc
// @Summary Search Steam games
// @Description Search for games on Steam by name
// @Tags games
// @Accept  json
// @Produce  json
// @Param q query string true "Search query"
// @Success 200 {object} model.APIResponse[[]model.SteamAppSearchResult]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /games/steam/search [get]
func (h *Handler) SearchSteamGames(c fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return c.Status(http.StatusBadRequest).JSON(model.Error{Error: "query parameter 'q' is required"})
	}

	encodedQuery := url.QueryEscape(query)
	steamURL := fmt.Sprintf("https://steamcommunity.com/actions/SearchApps/%s", encodedQuery)

	resp, err := http.Get(steamURL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to search steam: " + err.Error()})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to read steam response: " + err.Error()})
	}

	var steamResults []model.SteamAppSearchResult
	if err := json.Unmarshal(body, &steamResults); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: "failed to parse steam response: " + err.Error()})
	}

	return c.JSON(Response(steamResults, nil))
}

// CreateGame godoc
// @Summary Create a new game
// @Description Create a new game from Steam data
// @Tags games
// @Accept  json
// @Produce  json
// @Param game body model.CreateGameRequest true "Game data"
// @Success 201 {object} model.APIResponse[model.GameResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /games [post]
func (h *Handler) CreateGame(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var requestBody model.CreateGameRequest
	if err := c.Bind().JSON(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ErrInvalidRequestBody)
	}

	// Check if game with same ID already exists
	existingGame, err := h.repo.GetGameByID(ctx, requestBody.ID)
	if err == nil && existingGame != nil {
		return c.JSON(Response(existingGame, nil))
	}

	// Create new game
	id, err := h.repo.CreateGame(ctx, requestBody)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	game, err := h.repo.GetGameByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(model.Error{Error: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(Response(game, nil))
}
