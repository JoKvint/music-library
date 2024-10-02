package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"music-library/database"
	"music-library/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddSongRequest struct for adding a new song
type AddSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

// RegisterRoutes registers the API routes
func RegisterRoutes(router *gin.Engine) {
	router.GET("/songs", GetSongs)
	router.GET("/songs/:id", GetSongByID)
	router.POST("/songs", AddSong)
	router.PUT("/songs/:id", UpdateSong)
	router.DELETE("/songs/:id", DeleteSong)
}

// GetSongs retrieves a list of songs with optional filters
// @Summary Get list of songs
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group"
// @Param song query string false "Song"
// @Success 200 {array} models.Song
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	group := c.Query("group")
	title := c.Query("song")

	var songs []models.Song
	db := database.GetDB()

	if group != "" {
		db = db.Where("group = ?", group)
	}
	if title != "" {
		db = db.Where("title = ?", title)
	}

	if err := db.Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// GetSongByID retrieves a single song by its ID
// @Summary Get song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [get]
func GetSongByID(c *gin.Context) {
	id := c.Param("id")
	var song models.Song

	if err := database.GetDB().First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// AddSong adds a new song to the database
// @Summary Add a new song
// @Tags songs
// @Accept json
// @Produce json
// @Param song body AddSongRequest true "Add Song"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var request AddSongRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Запрос в внешний API
	url := fmt.Sprintf("http://external-api.com/info?group=%s&song=%s", request.Group, request.Song)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch song details")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details"})
		return
	}
	defer resp.Body.Close()

	var songDetail models.Song
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &songDetail); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse song details"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get song details"})
		return
	}

	// Сохранение в базе данных
	newSong := models.Song{
		Group:       request.Group,
		Title:       request.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}
	database.GetDB().Create(&newSong)

	c.JSON(http.StatusOK, newSong)
}

// UpdateSong updates a song's information
// @Summary Update a song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Update Song"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song

	if err := database.GetDB().First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	database.GetDB().Save(&song)
	c.JSON(http.StatusOK, song)
}

// DeleteSong deletes a song from the database
// @Summary Delete a song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song

	if err := database.GetDB().First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	database.GetDB().Delete(&song)
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
