package repository

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/sqlite"

	"github.com/K0ng2/zeedzad/model"
	repoModel "github.com/K0ng2/zeedzad/repository/model"
	. "github.com/K0ng2/zeedzad/repository/table"
)

type VideoWithGame struct {
	repoModel.Videos
	Game *repoModel.Games `alias:"game"`
}

func selectVideos() sqlite.SelectStatement {
	return sqlite.SELECT(
		Videos.AllColumns,
		Games.ID.AS("game.id"),
		Games.Name.AS("game.name"),
		Games.URL.AS("game.url"),
	).FROM(
		Videos.LEFT_JOIN(Games, Games.ID.EQ(Videos.GameID)),
	)
}

func (r *Repository) GetVideos(ctx context.Context, query model.Offset, search string) ([]model.VideoResponse, error) {
	var videos []VideoWithGame

	stmt := selectVideos()

	if search != "" {
		searchPattern := sqlite.String("%" + search + "%")
		stmt = stmt.WHERE(
			sqlite.OR(
				ILIKE(Videos.Title, searchPattern),
				ILIKE(Games.Name, searchPattern),
			),
		)
	}

	stmt = stmt.
		ORDER_BY(Videos.PublishedAt.DESC()).
		LIMIT(query.Limit).
		OFFSET(query.Offset)

	err := stmt.QueryContext(ctx, r.ex, &videos)
	if err != nil {
		return nil, FormatError("get videos", err)
	}

	return convertToVideoResponses(videos), nil
}

func (r *Repository) GetVideoByID(ctx context.Context, id string) (*model.VideoResponse, error) {
	var video VideoWithGame

	stmt := selectVideos().WHERE(Videos.ID.EQ(sqlite.String(id)))

	err := stmt.QueryContext(ctx, r.ex, &video)
	if err != nil {
		return nil, FormatError("get video by id", err)
	}

	responses := convertToVideoResponses([]VideoWithGame{video})
	if len(responses) == 0 {
		return nil, FormatError("get video by id", err)
	}

	return &responses[0], nil
}

func (r *Repository) GetVideoTotalItems(ctx context.Context, search string) (int64, error) {
	var expression *sqlite.BoolExpression

	if search != "" {
		searchPattern := sqlite.String("%" + search + "%")
		exp := ILIKE(Videos.Title, searchPattern)
		expression = &exp
	}

	return TotalItems(ctx, r.ex, Videos.ID, Videos, expression)
}

func (r *Repository) UpdateVideoGame(ctx context.Context, videoID string, gameID int64) error {
	stmt := Videos.UPDATE(Videos.GameID, Videos.UpdatedAt).
		SET(
			sqlite.Int(gameID),
			sqlite.CURRENT_TIMESTAMP(),
		).
		WHERE(Videos.ID.EQ(sqlite.String(videoID)))

	_, err := stmt.ExecContext(ctx, r.ex)
	if err != nil {
		return FormatError("update video game", err)
	}

	return nil
}

func (r *Repository) CreateVideo(ctx context.Context, video repoModel.Videos) error {
	stmt := Videos.INSERT(Videos.ID, Videos.Title, Videos.Thumbnail, Videos.PublishedAt, Videos.GameID, Videos.CreatedAt, Videos.UpdatedAt).
		VALUES(
			video.ID,
			video.Title,
			video.Thumbnail,
			video.PublishedAt,
			video.GameID,
			time.Now(),
			time.Now(),
		)

	_, err := stmt.ExecContext(ctx, r.ex)
	if err != nil {
		return FormatError("create video", err)
	}

	return nil
}

func (r *Repository) GetVideoByYouTubeID(ctx context.Context, youtubeID string) (*model.VideoResponse, error) {
	var video VideoWithGame

	stmt := selectVideos().WHERE(Videos.ID.EQ(sqlite.String(youtubeID)))

	err := stmt.QueryContext(ctx, r.ex, &video)
	if err != nil {
		return nil, FormatError("get video by youtube id", err)
	}

	responses := convertToVideoResponses([]VideoWithGame{video})
	if len(responses) == 0 {
		return nil, nil
	}

	return &responses[0], nil
}

func convertToVideoResponses(videos []VideoWithGame) []model.VideoResponse {
	responses := make([]model.VideoResponse, 0, len(videos))

	for _, v := range videos {
		response := model.VideoResponse{
			ID:          *v.ID,
			Title:       v.Title,
			Thumbnail:   v.Thumbnail,
			PublishedAt: v.PublishedAt,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}

		if v.Game != nil && v.Game.ID != nil {
			response.Game = &model.GameInfo{
				ID:   *v.Game.ID,
				Name: v.Game.Name,
				URL:  &v.Game.URL,
			}
		}

		responses = append(responses, response)
	}

	return responses
}
