package handlers

import (
	"net/http"
	"time"

	"kang-hyun-ji-backend/db"
	"kang-hyun-ji-backend/models"

	"github.com/gin-gonic/gin"
)

// GetMe godoc
// @Summary     내 정보 조회
// @Description 테스트 유저(id=1)의 정보를 반환합니다.
// @Tags        users
// @Produce     json
// @Success     200  {object}  models.User
// @Failure     404  {object}  map[string]string
// @Router      /users/me [get]
func GetMe(c *gin.Context) {
	user, err := fetchUser(1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateDepth godoc
// @Summary     수심 변경
// @Description 유저의 현재 수심을 변경합니다. (0~1000m)
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       body  body      map[string]int  true  "depth (0~1000)"
// @Success     200   {object}  models.User
// @Failure     400   {object}  map[string]string
// @Router      /users/me/depth [patch]
func UpdateDepth(c *gin.Context) {
	var req struct {
		Depth int `json:"depth" binding:"required,min=0,max=1000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(
		`UPDATE users SET current_depth = ?, updated_at = ? WHERE id = 1`,
		req.Depth, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check depth-related achievements
	checkDepthAchievements(1, req.Depth)

	user, _ := fetchUser(1)
	c.JSON(http.StatusOK, user)
}

func fetchUser(id int64) (*models.User, error) {
	row := db.DB.QueryRow(`
		SELECT id, username, current_depth, total_diary_count,
		       current_streak, longest_streak, last_diary_date,
		       created_at, updated_at
		FROM users WHERE id = ?`, id)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Username, &u.CurrentDepth, &u.TotalDiaryCount,
		&u.CurrentStreak, &u.LongestStreak, &u.LastDiaryDate,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
