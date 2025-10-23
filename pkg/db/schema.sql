-- Create tables for video and game matching application

-- Games table - stores Steam game information
CREATE TABLE IF NOT EXISTS games (
	id TEXT PRIMARY KEY,
	app_id TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	icon TEXT,
	logo TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_games_name ON games(name);
CREATE INDEX IF NOT EXISTS idx_games_app_id ON games(app_id);

-- Videos table - stores YouTube video information
CREATE TABLE IF NOT EXISTS videos (
	id TEXT PRIMARY KEY,
	youtube_id TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	description TEXT,
	thumbnail TEXT,
	published_at DATETIME NOT NULL,
	channel_id TEXT NOT NULL,
	channel_title TEXT,
	game_id TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_videos_youtube_id ON videos(youtube_id);
CREATE INDEX IF NOT EXISTS idx_videos_game_id ON videos(game_id);
CREATE INDEX IF NOT EXISTS idx_videos_published_at ON videos(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_title ON videos(title);
