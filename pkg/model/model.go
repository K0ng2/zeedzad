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
