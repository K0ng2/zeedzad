package model

import "time"

type INT64 struct {
	Number int64
}

type Meta struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type Offset struct {
	Limit  int64 `query:"limit,default:20"`
	Offset int64 `query:"offset,default:0"`
}

type APIResponse[T any] struct {
	Data T     `json:"data"`
	Meta *Meta `json:"meta,omitempty"` // omitted if nil
}

type DatabaseHealth struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

type Error struct {
	Error string `json:"error"`
}

// Video related models
type VideoResponse struct {
	ID           string    `json:"id"`
	YoutubeID    string    `json:"youtube_id"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	Thumbnail    *string   `json:"thumbnail"`
	PublishedAt  time.Time `json:"published_at"`
	ChannelID    string    `json:"channel_id"`
	ChannelTitle *string   `json:"channel_title"`
	Game         *GameInfo `json:"game"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GameInfo struct {
	ID    string  `json:"id"`
	AppID string  `json:"app_id"`
	Name  string  `json:"name"`
	Icon  *string `json:"icon"`
	Logo  *string `json:"logo"`
}

type UpdateVideoGameRequest struct {
	GameID string `json:"game_id"`
}

// Game related models
type GameResponse struct {
	ID        string    `json:"id"`
	AppID     string    `json:"app_id"`
	Name      string    `json:"name"`
	Icon      *string   `json:"icon"`
	Logo      *string   `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateGameRequest struct {
	AppID string  `json:"app_id"`
	Name  string  `json:"name"`
	Icon  *string `json:"icon"`
	Logo  *string `json:"logo"`
}

// Steam API response
type SteamAppSearchResult struct {
	AppID string `json:"appid"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Logo  string `json:"logo"`
}

// Sync result
type SyncResult struct {
	Added   int `json:"added"`
	Skipped int `json:"skipped"`
	Errors  int `json:"errors"`
	Total   int `json:"total"`
}
