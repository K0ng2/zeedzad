package server

import (
	"io/fs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/K0ng2/zeedzad/docs"
	"github.com/K0ng2/zeedzad/handler"
	"github.com/K0ng2/zeedzad/web"
)

// @BasePath /api
func NewRouter(handler *handler.Handler) *fiber.App {
	app := fiber.New()

	// Set up the Fiber app with middlewares
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New())
	app.Use(recover.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())

	api := app.Group("api")

	api.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler(
		httpSwagger.InstanceName(docs.SwaggerInfo.InfoInstanceName),
		httpSwagger.DefaultModelsExpandDepth(-1),
	)))

	api.Get("databasez", handler.DatabaseHealth)

	// Video routes
	api.Get("/videos", handler.GetVideos)
	api.Post("/videos/sync", handler.SyncYouTubeVideos)
	api.Get("/videos/:id", handler.GetVideoByID)
	api.Put("/videos/:id/game", handler.UpdateVideoGame)

	// Game routes
	api.Get("/games", handler.GetGames)
	api.Get("/games/:id", handler.GetGameByID)
	api.Post("/games", handler.CreateGame)
	api.Get("/games/steam/search", handler.SearchSteamGames)

	fSys, err := fs.Sub(web.EmbeddedFiles, "public")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}

	app.Use(static.New("", static.Config{
		FS:       fSys,
		Compress: true,
		NotFoundHandler: func(c fiber.Ctx) error {
			return c.SendFile("404.html", fiber.SendFile{
				FS:       fSys,
				Compress: true,
			})
		},
	}))

	return app
}
