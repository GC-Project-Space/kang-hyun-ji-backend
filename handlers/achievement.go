package handlers

import (
	"net/http"

	"kang-hyun-ji-backend/db"
	"kang-hyun-ji-backend/models"

	"github.com/gin-gonic/gin"
)

// ListAchievements returns all achievements with the user's completion status
func ListAchievements(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT a.id, a.key, a.name, a.description, a.condition_type, a.condition_value, COALESCE(a.icon_url, ''),
		       ua.achieved_at
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.id = ua.achievement_id AND ua.user_id = 1
		ORDER BY a.id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type AchievementWithStatus struct {
		models.Achievement
		Achieved   bool    `json:"achieved"`
		AchievedAt *string `json:"achieved_at"`
	}

	result := []AchievementWithStatus{}
	for rows.Next() {
		var entry AchievementWithStatus
		var achievedAt *string
		rows.Scan(
			&entry.ID, &entry.Key, &entry.Name, &entry.Description,
			&entry.ConditionType, &entry.ConditionValue, &entry.IconURL,
			&achievedAt,
		)
		entry.AchievedAt = achievedAt
		entry.Achieved = achievedAt != nil
		result = append(result, entry)
	}
	c.JSON(http.StatusOK, result)
}

// checkCollectionAchievements checks and grants collection-based achievements
func checkCollectionAchievements(userID int64) {
	// Total collected (unique creatures)
	var uniqueCount int
	db.DB.QueryRow(`
		SELECT COUNT(DISTINCT creature_id) FROM user_collection WHERE user_id = ?
	`, userID).Scan(&uniqueCount)

	grantIfNotExists(userID, "first_collect", uniqueCount >= 1)
	grantIfNotExists(userID, "collect_5", uniqueCount >= 5)
	grantIfNotExists(userID, "collect_10", uniqueCount >= 10)
	grantIfNotExists(userID, "collect_30", uniqueCount >= 30)

	// Check full collection
	var totalCreatures int
	db.DB.QueryRow(`SELECT COUNT(*) FROM sea_creatures`).Scan(&totalCreatures)
	grantIfNotExists(userID, "full_collection", uniqueCount >= totalCreatures)

	// Zone completions
	checkZoneCompletion(userID, "zone_shallow", 0, 250)
	checkZoneCompletion(userID, "zone_midwater", 250, 500)
	checkZoneCompletion(userID, "zone_deep", 500, 750)
	checkZoneCompletion(userID, "zone_abyss", 750, 1000)
}

func checkZoneCompletion(userID int64, key string, minDepth, maxDepth int) {
	var totalInZone int
	db.DB.QueryRow(`
		SELECT COUNT(*) FROM sea_creatures
		WHERE min_depth >= ? AND max_depth <= ?
	`, minDepth, maxDepth).Scan(&totalInZone)

	if totalInZone == 0 {
		// Creatures that span zone boundaries — count creatures primarily in this zone
		db.DB.QueryRow(`
			SELECT COUNT(*) FROM sea_creatures
			WHERE (min_depth + max_depth) / 2 >= ? AND (min_depth + max_depth) / 2 < ?
		`, minDepth, maxDepth).Scan(&totalInZone)
	}

	if totalInZone == 0 {
		return
	}

	var collectedInZone int
	db.DB.QueryRow(`
		SELECT COUNT(DISTINCT sc.id)
		FROM user_collection uc
		JOIN sea_creatures sc ON uc.creature_id = sc.id
		WHERE uc.user_id = ? AND (sc.min_depth + sc.max_depth) / 2 >= ? AND (sc.min_depth + sc.max_depth) / 2 < ?
	`, userID, minDepth, maxDepth).Scan(&collectedInZone)

	grantIfNotExists(userID, key, collectedInZone >= totalInZone)
}

func checkStreakAchievements(userID int64, streak int) {
	grantIfNotExists(userID, "streak_3", streak >= 3)
	grantIfNotExists(userID, "streak_7", streak >= 7)
	grantIfNotExists(userID, "streak_30", streak >= 30)
}

func checkDiaryCountAchievements(userID int64, count int) {
	grantIfNotExists(userID, "diary_10", count >= 10)
	grantIfNotExists(userID, "diary_50", count >= 50)
}

func checkDepthAchievements(userID int64, depth int) {
	grantIfNotExists(userID, "depth_250", depth >= 250)
	grantIfNotExists(userID, "depth_500", depth >= 500)
	grantIfNotExists(userID, "depth_750", depth >= 750)
	grantIfNotExists(userID, "depth_1000", depth >= 1000)
}

func grantIfNotExists(userID int64, key string, condition bool) {
	if !condition {
		return
	}
	var achievementID int64
	err := db.DB.QueryRow(`SELECT id FROM achievements WHERE key = ?`, key).Scan(&achievementID)
	if err != nil {
		return
	}
	db.DB.Exec(`
		INSERT OR IGNORE INTO user_achievements (user_id, achievement_id)
		VALUES (?, ?)
	`, userID, achievementID)
}
