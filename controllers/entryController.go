package controllers

import (
	"fmt"
	"go-money-tracker/initializers"
	"go-money-tracker/middleware"
	"go-money-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func extractEid(c *gin.Context) int {
	// Get entry id from route param
	_eid := c.Param("eid")

	if _eid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find entry id",
		})
		return 0
	}

	eid, err := strconv.Atoi(_eid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse entry id",
		})
		return 0
	}

	return eid
}

func CreateEntry(c *gin.Context) {
	// Parse Uid from headers
	uid, err := strconv.Atoi(c.Request.Header["Uid"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Uid not found in headers",
		})
		return
	}

	// Get entry body
	var body struct {
		Operation  string
		Amount     float64
		Source     string
		Target     string
		Category   string
		Note       string
		Datestring string
		Timestamp  int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Create new entry
	entry := models.Entry{
		UID:        uid,
		Operation:  body.Operation,
		Amount:     body.Amount,
		Source:     body.Source,
		Target:     body.Target,
		Category:   body.Category,
		Note:       body.Note,
		Datestring: body.Datestring,
		Timestamp:  body.Timestamp,
	}
	result := initializers.DB.Create(&entry)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create entry",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Entry created",
	})
}

func GetEntry(c *gin.Context) {
	// Extract UID from req headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		return
	}

	// Get entry id from route param
	eid := extractEid(c)
	if eid == 0 {
		return
	}

	// Look up target entry
	var entry models.Entry
	initializers.DB.Where(&models.Entry{ID: eid}).First(&entry)

	if eid == 0 || entry.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Can't find this entry",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"entry": entry,
	})
}

func GetEntries(c *gin.Context) {
	// Extract UID from req headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		return
	}

	// Get request body
	var body struct {
		DateStart  int    `json:"date_start"`
		DateEnd    int    `json:"date_end"`
		Datestring string `json:"datestring"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up matching entries
	entries := []models.Entry{}
	if body.Datestring != "" {
		datestring := fmt.Sprintf("%s%s%s", "%", body.Datestring, "%")
		initializers.DB.Where("datestring LIKE ?", datestring).Find(&entries, "uid = ?", uid)
	} else {
		initializers.DB.Where("timestamp >= ? AND timestamp <= ?", body.DateStart, body.DateEnd).Find(&entries, "uid = ?", uid)
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
	})
}

func UpdateEntry(c *gin.Context) {
	// Get entry id from router param
	eid := extractEid(c)
	if eid == 0 {
		return
	}

	// Get entry body
	var body struct {
		Operation  string
		Amount     float64
		Source     string
		Target     string
		Category   string
		Note       string
		Datestring string
		Timestamp  int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up target entry
	var entry models.Entry
	initializers.DB.First(&entry, "id = ?", eid)

	if entry.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Can't find this entry",
		})
		return
	}

	// Update entry
	result := initializers.DB.Model(&entry).Updates(models.Entry{
		Operation:  body.Operation,
		Amount:     body.Amount,
		Source:     body.Source,
		Target:     body.Target,
		Category:   body.Category,
		Note:       body.Note,
		Datestring: body.Datestring,
		Timestamp:  body.Timestamp,
	})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update entry",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Entry updated",
	})
}

func DeleteEntry(c *gin.Context) {
	// Get entry id from route param
	id := extractEid(c)
	if id == 0 {
		return
	}

	// Look up target entry
	var entry models.Entry
	initializers.DB.First(&entry, "id = ?", id)

	if entry.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Can't find this entry",
		})
		return
	}

	// Delete entry
	result := initializers.DB.Delete(&entry)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete entry",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Entry deleted",
	})
}

func GetAssets(c *gin.Context) {
	// Extract Uid from request headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user with this uid",
		})
		return
	}

	// Get assets (target where operation is income OR source where operation is expense OR source or target where operation is transfer) with unique categories related to Uid
	var incomeEntries []models.Entry
	initializers.DB.Where("operation = ?", "income").Distinct("target").Find(&incomeEntries, "uid = ?", uid)
	var expenseEntries []models.Entry
	initializers.DB.Where("operation = ?", "expense").Distinct("source").Find(&expenseEntries, "uid = ?", uid)
	var transferEntries []models.Entry
	initializers.DB.Where("operation = ?", "expense").Distinct("source", "target").Find(&transferEntries, "uid = ?", uid)

	// Map each category into a slice
	assets := []string{"account", "cash"}
	// income > target
	for _, entry := range incomeEntries {
		if entry.Target != "" && !slices.Contains(assets, entry.Target) {
			assets = append(assets, entry.Target)
		}
	}
	// expense > source
	for _, entry := range expenseEntries {
		if entry.Source != "" && !slices.Contains(assets, entry.Source) {
			assets = append(assets, entry.Source)
		}
	}
	// transfer > source OR target
	for _, entry := range transferEntries {
		if entry.Source != "" && !slices.Contains(assets, entry.Source) {
			assets = append(assets, entry.Source)
		}
		if entry.Target != "" && !slices.Contains(assets, entry.Target) {
			assets = append(assets, entry.Target)
		}
	}
	if !slices.Contains(assets, "other") {
		assets = append(assets, "other")
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"assets": assets,
	})
}

func GetNotes(c *gin.Context) {
	// Extract Uid from request headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user with this uid",
		})
		return
	}

	// Get entries with unique notes related to Uid
	var entries []models.Entry
	initializers.DB.Distinct("note").Find(&entries, "uid = ?", uid)

	// Map each note into a slice
	notes := []string{}
	for _, entry := range entries {
		if entry.Note != "" {
			notes = append(notes, entry.Note)
		}
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func GetExpenseCategories(c *gin.Context) {
	// Extract Uid from request headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user with this uid",
		})
		return
	}

	// Get entries where operation is expense with unique category related to Uid
	var entries []models.Entry
	initializers.DB.Where("operation = ?", "expense").Distinct("category").Find(&entries, "uid = ?", uid)

	// Map each note into a slice
	categories := []string{"education", "food", "gift", "health", "household", "social", "transport"}
	for _, entry := range entries {
		if entry.Category != "" && !slices.Contains(categories, entry.Category) {
			categories = append(categories, entry.Category)
		}
	}
	if !slices.Contains(categories, "other") {
		categories = append(categories, "other")
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func GetIncomeCategories(c *gin.Context) {
	// Extract Uid from request headers
	uid := middleware.ExtractUid(c)
	if uid == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user with this uid",
		})
		return
	}

	// Get entries where operation is income with unique category related to Uid
	var entries []models.Entry
	initializers.DB.Where("operation = ?", "income").Distinct("category").Find(&entries, "uid = ?", uid)

	// Map each note into a slice
	categories := []string{"bonus", "commision", "petty cash", "salary"}
	for _, entry := range entries {
		if entry.Category != "" && !slices.Contains(categories, entry.Category) {
			categories = append(categories, entry.Category)
		}
	}
	if !slices.Contains(categories, "other") {
		categories = append(categories, "other")
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
