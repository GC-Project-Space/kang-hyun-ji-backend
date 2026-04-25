package handlers

import (
	"net/http"

	"kang-hyun-ji-backend/db"
	"kang-hyun-ji-backend/models"

	"github.com/gin-gonic/gin"
)

// ListCreatures returns all creatures in the guidebook (도감)
func ListCreatures(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, name, description, min_depth, max_depth, rarity, spawn_weight, image_url
		FROM sea_creatures
		ORDER BY min_depth, max_depth
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	creatures := []models.SeaCreature{}
	for rows.Next() {
		var sc models.SeaCreature
		rows.Scan(&sc.ID, &sc.Name, &sc.Description, &sc.MinDepth, &sc.MaxDepth,
			&sc.Rarity, &sc.SpawnWeight, &sc.ImageURL)
		creatures = append(creatures, sc)
	}
	c.JSON(http.StatusOK, creatures)
}

// GetCollection returns the user's collected creatures grouped by creature with counts
func GetCollection(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT sc.id, sc.name, sc.description, sc.min_depth, sc.max_depth,
		       sc.rarity, sc.spawn_weight, sc.image_url,
		       COUNT(uc.id) as count,
		       MAX(uc.collected_at) as last_collected_at
		FROM sea_creatures sc
		LEFT JOIN user_collection uc ON sc.id = uc.creature_id AND uc.user_id = 1
		GROUP BY sc.id
		ORDER BY sc.min_depth, sc.max_depth
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type CollectionEntry struct {
		models.SeaCreature
		Count           int     `json:"count"`
		LastCollectedAt *string `json:"last_collected_at"`
		Collected       bool    `json:"collected"`
	}

	result := []CollectionEntry{}
	for rows.Next() {
		var entry CollectionEntry
		var lastCollectedAt *string
		rows.Scan(
			&entry.ID, &entry.Name, &entry.Description, &entry.MinDepth, &entry.MaxDepth,
			&entry.Rarity, &entry.SpawnWeight, &entry.ImageURL,
			&entry.Count, &lastCollectedAt,
		)
		entry.LastCollectedAt = lastCollectedAt
		entry.Collected = entry.Count > 0
		result = append(result, entry)
	}
	c.JSON(http.StatusOK, result)
}
