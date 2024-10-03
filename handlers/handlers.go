package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"music-library/database"
	"music-library/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

func RegisterRoutes(router *gin.Engine) {
	router.GET("/songs", GetSongs)
	router.GET("/songs/:id", GetSongByID)
	router.POST("/songs", AddSong)
	router.PUT("/songs/:id", UpdateSong)
	router.DELETE("/songs/:id", DeleteSong)
}

func GetSongs(c *gin.Context) {
	group := c.Query("group")
	title := c.Query("song")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var songs []models.Song
	db := database.GetDB()

	if group != "" {
		db = db.Where("group = ?", group)
	}
	if title != "" {
		db = db.Where("title = ?", title)
	}

	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func GetSongByID(c *gin.Context) {
	id := c.Param("id")
	var song models.Song

	if err := database.GetDB().First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

func AddSong(c *gin.Context) {
	var request AddSongRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

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
