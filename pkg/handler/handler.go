package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/K0ng2/zeedzad/db"
	"github.com/K0ng2/zeedzad/igdb"
	"github.com/K0ng2/zeedzad/model"
	"github.com/K0ng2/zeedzad/repository"
)

type Handler struct {
	repo *repository.Repository
	igdb *igdb.Client
}

func NewHandler(db *db.Database, igdbClient *igdb.Client) *Handler {
	return &Handler{
		repo: repository.NewRepository(db),
		igdb: igdbClient,
	}
}

// GetOffset extracts limit and offset query parameters from the request context.
func GetOffset(c fiber.Ctx) (*model.Offset, error) {
	var query model.Offset

	if err := c.Bind().Query(&query); err != nil {
		return nil, err
	}

	return &query, nil
}

var startTime = time.Now()

func DatabaseHealth(ping error, c fiber.Ctx) error {
	// Check database connectivity
	dbStatus := "healthy"
	if err := ping; err != nil {
		dbStatus = "unhealthy"
	}

	// Calculate uptime
	uptime := time.Since(startTime).String()

	response := model.DatabaseHealth{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  dbStatus,
		Uptime:    uptime,
	}

	// If database is unhealthy, set overall status to unhealthy
	if dbStatus == "unhealthy" {
		response.Status = "degraded"
		return c.Status(http.StatusServiceUnavailable).JSON(response)
	}

	return c.JSON(response)
}

func Response[T any](r T, m *model.Meta) model.APIResponse[T] {
	return model.APIResponse[T]{
		Data: r,
		Meta: m,
	}
}

var (
	ErrInvalidPathParams  = model.Error{Error: "invalid path parameters"}
	ErrInvalidRequestBody = model.Error{Error: "invalid request body"}
	ErrInvalidQueryParams = model.Error{Error: "invalid query parameters"}
)
