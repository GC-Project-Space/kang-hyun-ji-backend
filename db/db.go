package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(path string) {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	DB.SetMaxOpenConns(1)

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	migrate()
	log.Println("database initialized")
}

func migrate() {
	schema := `
	PRAGMA foreign_keys = ON;
	PRAGMA journal_mode = WAL;

	CREATE TABLE IF NOT EXISTS users (
		id                INTEGER PRIMARY KEY AUTOINCREMENT,
		username          TEXT NOT NULL UNIQUE,
		current_depth     INTEGER NOT NULL DEFAULT 0,
		total_diary_count INTEGER NOT NULL DEFAULT 0,
		current_streak    INTEGER NOT NULL DEFAULT 0,
		longest_streak    INTEGER NOT NULL DEFAULT 0,
		last_diary_date   TEXT,
		created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sea_creatures (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT NOT NULL,
		description  TEXT,
		min_depth    INTEGER NOT NULL,
		max_depth    INTEGER NOT NULL,
		rarity       TEXT NOT NULL DEFAULT 'COMMON',
		spawn_weight INTEGER NOT NULL DEFAULT 10,
		image_url    TEXT,
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS diaries (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id     INTEGER NOT NULL REFERENCES users(id),
		title       TEXT NOT NULL,
		content     TEXT NOT NULL,
		diary_date  TEXT NOT NULL,
		category    TEXT NOT NULL DEFAULT '일상',
		depth       INTEGER NOT NULL,
		creature_id INTEGER REFERENCES sea_creatures(id),
		hatches_at  DATETIME NOT NULL,
		is_hatched  INTEGER NOT NULL DEFAULT 0,
		hatched_at  DATETIME,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_collection (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id      INTEGER NOT NULL REFERENCES users(id),
		creature_id  INTEGER NOT NULL REFERENCES sea_creatures(id),
		diary_id     INTEGER NOT NULL REFERENCES diaries(id),
		collected_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS achievements (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		key             TEXT NOT NULL UNIQUE,
		name            TEXT NOT NULL,
		description     TEXT,
		condition_type  TEXT NOT NULL,
		condition_value INTEGER,
		icon_url        TEXT,
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_achievements (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id        INTEGER NOT NULL REFERENCES users(id),
		achievement_id INTEGER NOT NULL REFERENCES achievements(id),
		achieved_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, achievement_id)
	);
	`

	if _, err := DB.Exec(schema); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
