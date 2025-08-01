package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/uso/uso/internal/database"
	"github.com/uso/uso/internal/models"
)

type SearchHandler struct {
	db *database.DB
}

func NewSearchHandler(db *database.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) SearchCasts(c *gin.Context) {
	var params models.CastSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 50 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	// Build query
	query := `
		SELECT DISTINCT u.id, u.name, u.profile_image, 
		       cp.id as profile_id, cp.bio, cp.hourly_rate, cp.rank, cp.service_areas,
		       COALESCE(AVG(r.rating), 0) as rating, COUNT(DISTINCT r.id) as review_count
		FROM users u
		JOIN cast_profiles cp ON u.id = cp.user_id
		LEFT JOIN reviews r ON r.reviewed_id = u.id
		WHERE u.user_type = 'cast' AND cp.approval_status = 'approved'
	`
	
	args := []interface{}{}
	argCount := 1

	// Location filter
	if params.Location != "" {
		query += fmt.Sprintf(" AND $%d = ANY(cp.service_areas)", argCount)
		args = append(args, params.Location)
		argCount++
	}

	// Price range filter
	if params.MinPrice > 0 {
		query += fmt.Sprintf(" AND cp.hourly_rate >= $%d", argCount)
		args = append(args, params.MinPrice)
		argCount++
	}
	if params.MaxPrice > 0 {
		query += fmt.Sprintf(" AND cp.hourly_rate <= $%d", argCount)
		args = append(args, params.MaxPrice)
		argCount++
	}

	// Rank filter
	if params.Rank != "" {
		query += fmt.Sprintf(" AND cp.rank = $%d", argCount)
		args = append(args, params.Rank)
		argCount++
	}

	// Date/time availability filter
	if !params.Date.IsZero() && params.StartTime != "" {
		// Exclude casts with conflicting bookings
		query += fmt.Sprintf(` AND u.id NOT IN (
			SELECT cast_id FROM bookings 
			WHERE booking_date = $%d 
			AND status IN ('pending', 'accepted')
			AND ($%d::time >= start_time AND $%d::time < start_time + (duration_hours || ' hours')::interval)
		)`, argCount, argCount+1, argCount+1)
		args = append(args, params.Date, params.StartTime)
		argCount += 2
	}

	query += " GROUP BY u.id, u.name, u.profile_image, cp.id, cp.bio, cp.hourly_rate, cp.rank, cp.service_areas"
	query += " ORDER BY rating DESC, review_count DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, params.Limit, offset)

	// Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error searching casts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	casts := []gin.H{}
	for rows.Next() {
		var userID, profileID int
		var name string
		var profileImage sql.NullString
		var bio sql.NullString
		var hourlyRate float64
		var rank models.CastRank
		var serviceAreas pq.StringArray
		var rating float64
		var reviewCount int

		err := rows.Scan(&userID, &name, &profileImage, &profileID, &bio, 
			&hourlyRate, &rank, &serviceAreas, &rating, &reviewCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Get first gallery image
		var galleryImage sql.NullString
		h.db.QueryRow(`
			SELECT image_url FROM cast_gallery_images 
			WHERE cast_profile_id = $1 
			ORDER BY display_order LIMIT 1
		`, profileID).Scan(&galleryImage)

		cast := gin.H{
			"id":            userID,
			"name":          name,
			"profile_image": profileImage.String,
			"gallery_image": galleryImage.String,
			"bio":           bio.String,
			"hourly_rate":   hourlyRate,
			"rank":          rank,
			"service_areas": serviceAreas,
			"rating":        rating,
			"review_count":  reviewCount,
		}
		casts = append(casts, cast)
	}

	// Get total count
	countQuery := strings.Replace(query, "SELECT DISTINCT u.id, u.name, u.profile_image, cp.id as profile_id, cp.bio, cp.hourly_rate, cp.rank, cp.service_areas, COALESCE(AVG(r.rating), 0) as rating, COUNT(DISTINCT r.id) as review_count", "SELECT COUNT(DISTINCT u.id)", 1)
	countQuery = strings.Split(countQuery, " GROUP BY")[0]
	countQuery = strings.Split(countQuery, " ORDER BY")[0]
	countQuery = strings.Split(countQuery, " LIMIT")[0]
	
	var totalCount int
	countArgs := args[:len(args)-2] // Remove LIMIT and OFFSET args
	err = h.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("Error getting count: %v", err)
		totalCount = len(casts)
	}

	c.JSON(http.StatusOK, gin.H{
		"casts":       casts,
		"total":       totalCount,
		"page":        params.Page,
		"limit":       params.Limit,
		"total_pages": (totalCount + params.Limit - 1) / params.Limit,
	})
}

func (h *SearchHandler) GetCastProfile(c *gin.Context) {
	castID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cast ID"})
		return
	}

	// Get cast details
	var profile models.UserProfile
	var castProfile models.CastProfile
	
	err = h.db.QueryRow(`
		SELECT u.id, u.email, u.user_type, u.name, u.phone, u.birth_date, u.profile_image,
		       cp.id, cp.user_id, cp.bio, cp.hourly_rate, cp.rank, cp.service_areas, 
		       cp.approval_status, cp.approved_at,
		       COALESCE(AVG(r.rating), 0) as rating, COUNT(r.id) as review_count
		FROM users u
		JOIN cast_profiles cp ON u.id = cp.user_id
		LEFT JOIN reviews r ON r.reviewed_id = u.id
		WHERE u.id = $1 AND u.user_type = 'cast' AND cp.approval_status = 'approved'
		GROUP BY u.id, cp.id
	`, castID).Scan(
		&profile.ID, &profile.Email, &profile.UserType, &profile.Name,
		&profile.Phone, &profile.BirthDate, &profile.ProfileImage,
		&castProfile.ID, &castProfile.UserID, &castProfile.Bio,
		&castProfile.HourlyRate, &castProfile.Rank, &castProfile.ServiceAreas,
		&castProfile.ApprovalStatus, &castProfile.ApprovedAt,
		&profile.Rating, &profile.ReviewCount,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast not found"})
		return
	} else if err != nil {
		log.Printf("Error getting cast profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get gallery images
	rows, err := h.db.Query(`
		SELECT id, image_url, display_order
		FROM cast_gallery_images
		WHERE cast_profile_id = $1
		ORDER BY display_order
	`, castProfile.ID)
	
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var img models.GalleryImage
			if err := rows.Scan(&img.ID, &img.ImageURL, &img.DisplayOrder); err == nil {
				castProfile.GalleryImages = append(castProfile.GalleryImages, img)
			}
		}
	}

	profile.CastProfile = &castProfile

	// Get recent reviews
	reviewRows, err := h.db.Query(`
		SELECT r.id, r.rating, r.comment, r.created_at, u.name
		FROM reviews r
		JOIN users u ON r.reviewer_id = u.id
		WHERE r.reviewed_id = $1
		ORDER BY r.created_at DESC
		LIMIT 10
	`, castID)

	reviews := []gin.H{}
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var review models.Review
			var reviewerName string
			if err := reviewRows.Scan(&review.ID, &review.Rating, &review.Comment, 
				&review.CreatedAt, &reviewerName); err == nil {
				reviews = append(reviews, gin.H{
					"id":            review.ID,
					"rating":        review.Rating,
					"comment":       review.Comment,
					"created_at":    review.CreatedAt,
					"reviewer_name": reviewerName,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
		"reviews": reviews,
	})
}

func (h *SearchHandler) GetServiceAreas(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, name FROM service_areas ORDER BY display_order
	`)
	if err != nil {
		log.Printf("Error getting service areas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	areas := []models.ServiceArea{}
	for rows.Next() {
		var area models.ServiceArea
		if err := rows.Scan(&area.ID, &area.Name); err == nil {
			areas = append(areas, area)
		}
	}

	c.JSON(http.StatusOK, areas)
}