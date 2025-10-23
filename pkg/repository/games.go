package repository

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/sqlite"
	"github.com/google/uuid"

	"github.com/K0ng2/zeedzad/model"
	repoModel "github.com/K0ng2/zeedzad/repository/model"
	. "github.com/K0ng2/zeedzad/repository/table"
)

func selectGames() sqlite.SelectStatement {
	return sqlite.SELECT(Games.AllColumns).FROM(Games)
}

func (r *Repository) GetGames(ctx context.Context, query model.Offset, search string) ([]model.GameResponse, error) {
	var games []repoModel.Games

	stmt := selectGames()

	if search != "" {
		searchPattern := sqlite.String("%" + search + "%")
		stmt = stmt.WHERE(ILIKE(Games.Name, searchPattern))
	}

	stmt = stmt.
		ORDER_BY(Games.Name.ASC()).
		LIMIT(query.Limit).
		OFFSET(query.Offset)

	err := stmt.QueryContext(ctx, r.ex, &games)
	if err != nil {
		return nil, FormatError("get games", err)
	}

	return convertToGameResponses(games), nil
}

func (r *Repository) GetGameByID(ctx context.Context, id uuid.UUID) (*model.GameResponse, error) {
	var game repoModel.Games

	stmt := selectGames().WHERE(Games.ID.EQ(sqlite.String(id.String())))

	err := stmt.QueryContext(ctx, r.ex, &game)
	if err != nil {
		return nil, FormatError("get game by id", err)
	}

	responses := convertToGameResponses([]repoModel.Games{game})
	if len(responses) == 0 {
		return nil, FormatError("get game by id", err)
	}

	return &responses[0], nil
}

func (r *Repository) GetGameTotalItems(ctx context.Context, search string) (int64, error) {
	var expression *sqlite.BoolExpression

	if search != "" {
		searchPattern := sqlite.String("%" + search + "%")
		exp := ILIKE(Games.Name, searchPattern)
		expression = &exp
	}

	return TotalItems(ctx, r.ex, Games.ID, Games, expression)
}

func (r *Repository) CreateGame(ctx context.Context, req model.CreateGameRequest) (string, error) {
	id := uuid.New().String()

	stmt := Games.INSERT(Games.ID, Games.AppID, Games.Name, Games.Icon, Games.Logo, Games.CreatedAt, Games.UpdatedAt).
		VALUES(
			id,
			req.AppID,
			req.Name,
			req.Icon,
			req.Logo,
			time.Now(),
			time.Now(),
		)

	_, err := stmt.ExecContext(ctx, r.ex)
	if err != nil {
		return "", FormatError("create game", err)
	}

	return id, nil
}

func (r *Repository) SearchGamesByAppID(ctx context.Context, appID string) (*model.GameResponse, error) {
	var game repoModel.Games

	stmt := selectGames().WHERE(Games.AppID.EQ(sqlite.String(appID)))

	err := stmt.QueryContext(ctx, r.ex, &game)
	if err != nil {
		return nil, FormatError("search game by app id", err)
	}

	responses := convertToGameResponses([]repoModel.Games{game})
	if len(responses) == 0 {
		return nil, nil
	}

	return &responses[0], nil
}

func convertToGameResponses(games []repoModel.Games) []model.GameResponse {
	responses := make([]model.GameResponse, 0, len(games))

	for _, g := range games {
		if g.ID == nil {
			continue
		}

		responses = append(responses, model.GameResponse{
			ID:        *g.ID,
			AppID:     g.AppID,
			Name:      g.Name,
			Icon:      g.Icon,
			Logo:      g.Logo,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		})
	}

	return responses
}
