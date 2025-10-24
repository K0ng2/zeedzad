package repository

import (
	"context"
	"time"

	"github.com/go-jet/jet/v2/sqlite"

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

func (r *Repository) GetGameByID(ctx context.Context, id int64) (*model.GameResponse, error) {
	var game repoModel.Games

	stmt := selectGames().WHERE(Games.ID.EQ(sqlite.Int(id)))

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

func (r *Repository) CreateGame(ctx context.Context, req model.CreateGameRequest) (int64, error) {
	var iconValue sqlite.Expression = sqlite.NULL
	if req.Icon != nil {
		iconValue = sqlite.String(*req.Icon)
	}

	var logoValue sqlite.Expression = sqlite.NULL
	if req.Logo != nil {
		logoValue = sqlite.String(*req.Logo)
	}

	stmt := Games.INSERT(Games.ID, Games.Name, Games.Icon, Games.Logo, Games.CreatedAt, Games.UpdatedAt).
		VALUES(
			req.ID,
			req.Name,
			iconValue,
			logoValue,
			time.Now(),
			time.Now(),
		)

	_, err := stmt.ExecContext(ctx, r.ex)
	if err != nil {
		return 0, FormatError("create game", err)
	}

	return req.ID, nil
}

func convertToGameResponses(games []repoModel.Games) []model.GameResponse {
	responses := make([]model.GameResponse, 0, len(games))

	for _, g := range games {
		if g.ID == nil {
			continue
		}

		responses = append(responses, model.GameResponse{
			ID:        *g.ID,
			Name:      g.Name,
			Icon:      g.Icon,
			Logo:      g.Logo,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		})
	}

	return responses
}
