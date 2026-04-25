package models

import "time"

type User struct {
	ID             int64      `json:"id"`
	Username       string     `json:"username"`
	CurrentDepth   int        `json:"current_depth"`
	TotalDiaryCount int       `json:"total_diary_count"`
	CurrentStreak  int        `json:"current_streak"`
	LongestStreak  int        `json:"longest_streak"`
	LastDiaryDate  *string    `json:"last_diary_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type SeaCreature struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MinDepth    int     `json:"min_depth"`
	MaxDepth    int     `json:"max_depth"`
	Rarity      string  `json:"rarity"`
	SpawnWeight int     `json:"spawn_weight"`
	ImageURL    string  `json:"image_url"`
}

type Diary struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	DiaryDate  string     `json:"diary_date"`
	Category   string     `json:"category"`
	Depth      int        `json:"depth"`
	CreatureID *int64     `json:"creature_id"`
	HatchesAt  time.Time  `json:"hatches_at"`
	IsHatched  bool       `json:"is_hatched"`
	HatchedAt  *time.Time `json:"hatched_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Creature   *SeaCreature `json:"creature,omitempty"`
}

type UserCollection struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	CreatureID  int64        `json:"creature_id"`
	DiaryID     int64        `json:"diary_id"`
	CollectedAt time.Time    `json:"collected_at"`
	Creature    *SeaCreature `json:"creature,omitempty"`
	Count       int          `json:"count,omitempty"`
}

type Achievement struct {
	ID             int64  `json:"id"`
	Key            string `json:"key"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ConditionType  string `json:"condition_type"`
	ConditionValue int    `json:"condition_value"`
	IconURL        string `json:"icon_url"`
}

type UserAchievement struct {
	ID            int64       `json:"id"`
	UserID        int64       `json:"user_id"`
	AchievementID int64       `json:"achievement_id"`
	AchievedAt    time.Time   `json:"achieved_at"`
	Achievement   *Achievement `json:"achievement,omitempty"`
}
