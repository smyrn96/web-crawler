package services

import (
	"backend-crawler/database"
	"backend-crawler/models"
	"log"
	"time"
)

var jobQueue = make(chan uint, 100)

func EnqueueURL(id uint) {
	jobQueue <- id
}

func StartWorker() {
	go func() {
		for id := range jobQueue {
			go process(id)
		}
	}()
}

func process(id uint) {
	var url models.URL
	if err := database.DB.First(&url, id).Error; err != nil {
		log.Printf("URL ID %d not found: %v", id, err)
		return
	}

	// Update status to running
	database.DB.Model(&url).Update("status", "running")

	// Crawl and analyze
	result, err := CrawlURL(url.Address)
	if err != nil {
		database.DB.Model(&url).Updates(models.URL{
			Status: "error",
			Error:  err.Error(),
		})
		log.Printf("Failed to crawl URL %s: %v", url.Address, err)
		return
	}

	// Update database with result
	database.DB.Model(&url).Updates(map[string]interface{}{
		"status":         "done",
		"title":          result.Title,
		"html_version":   result.HTMLVersion,
		"has_login_form": result.HasLoginForm,
		"internal_links": result.InternalLinks,
		"external_links": result.ExternalLinks,
		"broken_links":   result.BrokenLinks,
		"error":          "",
		"updated_at":     time.Now(),
	})
}
