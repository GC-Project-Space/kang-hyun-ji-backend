package handlers

import (
	"database/sql"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"kang-hyun-ji-backend/db"
	"kang-hyun-ji-backend/models"

	"github.com/gin-gonic/gin"
)

var validCategories = map[string]bool{
	"감사": true, "일상": true, "감정": true, "목표": true, "기타": true,
}

// CreateDiaryRequest is the request body for creating a diary.
type CreateDiaryRequest struct {
	Title     string `json:"title" example:"오늘의 일기"`
	Content   string `json:"content" example:"오늘은 바다가 맑았다."`
	DiaryDate string `json:"diary_date" example:"2026-04-25"`
	Category  string `json:"category" example:"일상" enums:"감사,일상,감정,목표,기타"`
}

// ListDiaries godoc
// @Summary     일기 목록 조회
// @Description 테스트 유저의 전체 일기 목록을 반환합니다.
// @Tags        diaries
// @Produce     json
// @Success     200  {array}   models.Diary
// @Router      /diaries [get]
func ListDiaries(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT d.id, d.user_id, d.title, d.content, d.diary_date, d.category,
		       d.depth, d.creature_id, d.hatches_at, d.is_hatched, d.hatched_at,
		       d.created_at, d.updated_at,
		       sc.id, sc.name, sc.description, sc.min_depth, sc.max_depth,
		       sc.rarity, sc.spawn_weight, sc.image_url
		FROM diaries d
		LEFT JOIN sea_creatures sc ON d.creature_id = sc.id
		WHERE d.user_id = 1
		ORDER BY d.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	diaries := []models.Diary{}
	for rows.Next() {
		d, err := scanDiary(rows)
		if err != nil {
			continue
		}
		diaries = append(diaries, *d)
	}
	c.JSON(http.StatusOK, diaries)
}

// GetDiary godoc
// @Summary     일기 단건 조회
// @Description ID로 일기를 조회합니다.
// @Tags        diaries
// @Produce     json
// @Param       id   path      int  true  "Diary ID"
// @Success     200  {object}  models.Diary
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /diaries/{id} [get]
func GetDiary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	diary, err := fetchDiary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "diary not found"})
		return
	}
	c.JSON(http.StatusOK, diary)
}

// CreateDiary godoc
// @Summary     일기 작성 (알 생성)
// @Description 일기를 작성하면 현재 수심 기준으로 생명체 알이 배정됩니다.
// @Tags        diaries
// @Accept      json
// @Produce     json
// @Param       body  body      handlers.CreateDiaryRequest  true  "일기 내용"
// @Success     201   {object}  models.Diary
// @Failure     400   {object}  map[string]string
// @Router      /diaries [post]
func CreateDiary(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		DiaryDate string `json:"diary_date" binding:"required"`
		Category  string `json:"category" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validCategories[req.Category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	// Get user's current depth
	var currentDepth int
	db.DB.QueryRow(`SELECT current_depth FROM users WHERE id = 1`).Scan(&currentDepth)

	// Pick a random creature based on depth and spawn_weight
	creatureID := pickCreature(currentDepth)

	now := time.Now()
	hatchesAt := now.Add(10 * time.Second)

	result, err := db.DB.Exec(`
		INSERT INTO diaries (user_id, title, content, diary_date, category, depth, creature_id, hatches_at, created_at, updated_at)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Title, req.Content, req.DiaryDate, req.Category, currentDepth, creatureID, hatchesAt, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	diaryID, _ := result.LastInsertId()

	// Update user stats + streak
	updateUserStats(1, req.DiaryDate)

	diary, _ := fetchDiary(diaryID)
	c.JSON(http.StatusCreated, diary)
}

// HatchDiary godoc
// @Summary     알 부화
// @Description 작성 후 24시간이 지난 일기의 알을 부화시킵니다.
// @Tags        diaries
// @Produce     json
// @Param       id   path      int  true  "Diary ID"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Router      /diaries/{id}/hatch [post]
func HatchDiary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	diary, err := fetchDiary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "diary not found"})
		return
	}

	if diary.IsHatched {
		c.JSON(http.StatusOK, gin.H{"message": "already hatched", "diary": diary})
		return
	}

	if time.Now().Before(diary.HatchesAt) {
		remaining := time.Until(diary.HatchesAt)
		c.JSON(http.StatusOK, gin.H{
			"message":           "not ready yet",
			"remaining_seconds": int(remaining.Seconds()),
			"hatches_at":        diary.HatchesAt,
		})
		return
	}

	// Hatch!
	now := time.Now()
	db.DB.Exec(`UPDATE diaries SET is_hatched = 1, hatched_at = ?, updated_at = ? WHERE id = ?`, now, now, id)

	// Add to collection
	db.DB.Exec(`
		INSERT INTO user_collection (user_id, creature_id, diary_id, collected_at)
		VALUES (1, ?, ?, ?)
	`, diary.CreatureID, id, now)

	// Check collection achievements
	checkCollectionAchievements(1)

	diary, _ = fetchDiary(id)
	c.JSON(http.StatusOK, gin.H{"message": "hatched!", "diary": diary})
}

// pickCreature selects a creature by weighted random from creatures that include the given depth
func pickCreature(depth int) *int64 {
	rows, err := db.DB.Query(`
		SELECT id, spawn_weight FROM sea_creatures
		WHERE min_depth <= ? AND max_depth >= ?
	`, depth, depth)
	if err != nil {
		return nil
	}
	defer rows.Close()

	type entry struct {
		id     int64
		weight int
	}
	var pool []entry
	totalWeight := 0

	for rows.Next() {
		var e entry
		rows.Scan(&e.id, &e.weight)
		pool = append(pool, e)
		totalWeight += e.weight
	}

	if len(pool) == 0 {
		return nil
	}

	r := rand.Intn(totalWeight)
	for _, e := range pool {
		r -= e.weight
		if r < 0 {
			id := e.id
			return &id
		}
	}
	id := pool[len(pool)-1].id
	return &id
}

func updateUserStats(userID int64, diaryDate string) {
	var lastDiaryDate *string
	var currentStreak, longestStreak int
	db.DB.QueryRow(
		`SELECT last_diary_date, current_streak, longest_streak FROM users WHERE id = ?`, userID,
	).Scan(&lastDiaryDate, &currentStreak, &longestStreak)

	newStreak := 1
	if lastDiaryDate != nil {
		last, _ := time.Parse("2006-01-02", *lastDiaryDate)
		today, _ := time.Parse("2006-01-02", diaryDate)
		diff := today.Sub(last).Hours() / 24
		if diff == 1 {
			newStreak = currentStreak + 1
		} else if diff == 0 {
			newStreak = currentStreak // same day, no change
		}
	}

	newLongest := longestStreak
	if newStreak > longestStreak {
		newLongest = newStreak
	}

	db.DB.Exec(`
		UPDATE users
		SET total_diary_count = total_diary_count + 1,
		    current_streak = ?,
		    longest_streak = ?,
		    last_diary_date = ?,
		    updated_at = ?
		WHERE id = ?
	`, newStreak, newLongest, diaryDate, time.Now(), userID)

	// Check streak and diary count achievements
	checkStreakAchievements(userID, newStreak)
	var totalCount int
	db.DB.QueryRow(`SELECT total_diary_count FROM users WHERE id = ?`, userID).Scan(&totalCount)
	checkDiaryCountAchievements(userID, totalCount)
}

func fetchDiary(id int64) (*models.Diary, error) {
	row := db.DB.QueryRow(`
		SELECT d.id, d.user_id, d.title, d.content, d.diary_date, d.category,
		       d.depth, d.creature_id, d.hatches_at, d.is_hatched, d.hatched_at,
		       d.created_at, d.updated_at,
		       sc.id, sc.name, sc.description, sc.min_depth, sc.max_depth,
		       sc.rarity, sc.spawn_weight, sc.image_url
		FROM diaries d
		LEFT JOIN sea_creatures sc ON d.creature_id = sc.id
		WHERE d.id = ?
	`, id)

	return scanDiaryRow(row)
}

func scanDiary(rows *sql.Rows) (*models.Diary, error) {
	var d models.Diary
	var creature models.SeaCreature
	var creatureID sql.NullInt64
	var scID sql.NullInt64

	err := rows.Scan(
		&d.ID, &d.UserID, &d.Title, &d.Content, &d.DiaryDate, &d.Category,
		&d.Depth, &creatureID, &d.HatchesAt, &d.IsHatched, &d.HatchedAt,
		&d.CreatedAt, &d.UpdatedAt,
		&scID, &creature.Name, &creature.Description, &creature.MinDepth, &creature.MaxDepth,
		&creature.Rarity, &creature.SpawnWeight, &creature.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	if creatureID.Valid {
		id := creatureID.Int64
		d.CreatureID = &id
	}
	if scID.Valid {
		creature.ID = scID.Int64
		d.Creature = &creature
	}
	return &d, nil
}

func scanDiaryRow(row *sql.Row) (*models.Diary, error) {
	var d models.Diary
	var creature models.SeaCreature
	var creatureID sql.NullInt64
	var scID sql.NullInt64

	err := row.Scan(
		&d.ID, &d.UserID, &d.Title, &d.Content, &d.DiaryDate, &d.Category,
		&d.Depth, &creatureID, &d.HatchesAt, &d.IsHatched, &d.HatchedAt,
		&d.CreatedAt, &d.UpdatedAt,
		&scID, &creature.Name, &creature.Description, &creature.MinDepth, &creature.MaxDepth,
		&creature.Rarity, &creature.SpawnWeight, &creature.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	if creatureID.Valid {
		id := creatureID.Int64
		d.CreatureID = &id
	}
	if scID.Valid {
		creature.ID = scID.Int64
		d.Creature = &creature
	}
	return &d, nil
}
