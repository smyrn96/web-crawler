package handlers

import (
	"backend-crawler/database"
	"backend-crawler/models"
	"backend-crawler/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BrokenLink struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

func InitDB() {
	database.Connect()
	database.DB.AutoMigrate(&models.URL{})
}

// POST /api/urls
func AddURL(c *gin.Context) {
	var body struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&body); err != nil || body.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	url := models.URL{
		Address: body.URL,
		Status:  "queued",
	}
	database.DB.Create(&url)

	// Use services to queue job
	services.EnqueueURL(url.ID)

	c.JSON(http.StatusAccepted, url)
}

// GET /api/urls
func ListURLs(c *gin.Context) {
	var urls []models.URL
	database.DB.Order("created_at desc").Find(&urls)
	c.JSON(http.StatusOK, urls)
}

func GetURLDetails(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Decode broken_links field
	var brokenLinks []BrokenLink
	_ = json.Unmarshal([]byte(url.BrokenLinks), &brokenLinks)

	response := gin.H{
		"id":              url.ID,
		"address":         url.Address,
		"status":          url.Status,
		"title":           url.Title,
		"html_version":    url.HTMLVersion,
		"has_login_form":  url.HasLoginForm,
		"internal_links":  url.InternalLinks,
		"external_links":  url.ExternalLinks,
		"broken_links":    brokenLinks,
		"error":           url.Error,
		"created_at":      url.CreatedAt,
		"updated_at":      url.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}

// POST /api/urls/:id/reanalyze
func ReanalyzeURL(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Reset status and queue again
	database.DB.Model(&url).Updates(models.URL{Status: "queued", Error: ""})
	services.EnqueueURL(url.ID)

	c.JSON(http.StatusAccepted, gin.H{"message": "Reanalysis started"})
}

// DELETE /api/urls
func DeleteURLs(c *gin.Context) {
	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.BindJSON(&body); err != nil || len(body.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := database.DB.Delete(&models.URL{}, body.IDs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": result.RowsAffected})
}
