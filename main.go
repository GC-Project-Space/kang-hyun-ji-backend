// @title           Ocean Diary API
// @version         1.0
// @description     강현지의 바다 일기 백엔드 API입니다.
// @host            localhost:8080
// @BasePath        /api

package main

import (
	"log"

	_ "kang-hyun-ji-backend/docs"
	"kang-hyun-ji-backend/db"
	"kang-hyun-ji-backend/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.Init("./ocean_diary.db")
	db.Seed()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		// User (hardcoded test user id=1)
		api.GET("/users/me", handlers.GetMe)
		api.PATCH("/users/me/depth", handlers.UpdateDepth)

		// Diaries (eggs)
		api.GET("/diaries", handlers.ListDiaries)
		api.POST("/diaries", handlers.CreateDiary)
		api.GET("/diaries/:id", handlers.GetDiary)
		api.POST("/diaries/:id/hatch", handlers.HatchDiary)

		// Collection & Guidebook (도감)
		api.GET("/creatures", handlers.ListCreatures)
		api.GET("/collection", handlers.GetCollection)

		// Achievements
		api.GET("/achievements", handlers.ListAchievements)
	}

	log.Println("server starting on :8080")
	r.Run(":8080")
}
