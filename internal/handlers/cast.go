package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uso/uso/config"
	"github.com/uso/uso/internal/database"
	"github.com/uso/uso/internal/models"
)

type CastHandler struct {
	db  *database.DB
	cfg *config.Config
}

func NewCastHandler(db *database.DB, cfg *config.Config) *CastHandler {
	return &CastHandler{db: db, cfg: cfg}
}

func (h *CastHandler) UpdateCastProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req models.CastProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get cast profile ID
	var profileID int
	err := h.db.QueryRow("SELECT id FROM cast_profiles WHERE user_id = $1", userID).Scan(&profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast profile not found"})
		return
	}

	// Update profile
	query := `UPDATE cast_profiles SET `
	args := []interface{}{}
	argCount := 1

	if req.Bio != nil {
		query += `bio = $` + strconv.Itoa(argCount) + `, `
		args = append(args, *req.Bio)
		argCount++
	}

	if req.Rank != nil {
		hourlyRate := req.Rank.GetHourlyRate()
		query += `rank = $` + strconv.Itoa(argCount) + `, hourly_rate = $` + strconv.Itoa(argCount+1) + `, `
		args = append(args, *req.Rank, hourlyRate)
		argCount += 2
	}

	if len(req.ServiceAreas) > 0 {
		query += `service_areas = $` + strconv.Itoa(argCount) + `, `
		args = append(args, req.ServiceAreas)
		argCount++
	}

	query += `updated_at = CURRENT_TIMESTAMP WHERE id = $` + strconv.Itoa(argCount)
	args = append(args, profileID)

	_, err = h.db.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating cast profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *CastHandler) UploadGalleryImage(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	// Get cast profile
	var profileID int
	err := h.db.QueryRow("SELECT id FROM cast_profiles WHERE user_id = $1", userID).Scan(&profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast profile not found"})
		return
	}

	// Check gallery image count
	var imageCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM cast_gallery_images WHERE cast_profile_id = $1", profileID).Scan(&imageCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if imageCount >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 gallery images allowed"})
		return
	}

	// In production, you would handle file upload here and store in S3/Cloudinary
	imageURL := c.PostForm("image_url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL required"})
		return
	}

	// Insert gallery image
	var imageID int
	err = h.db.QueryRow(`
		INSERT INTO cast_gallery_images (cast_profile_id, image_url, display_order)
		VALUES ($1, $2, $3) RETURNING id
	`, profileID, imageURL, imageCount).Scan(&imageID)

	if err != nil {
		log.Printf("Error inserting gallery image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": imageID,
		"image_url": imageURL,
		"display_order": imageCount,
	})
}

func (h *CastHandler) DeleteGalleryImage(c *gin.Context) {
	userID := c.GetInt("user_id")
	imageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Verify ownership
	var ownerID int
	err = h.db.QueryRow(`
		SELECT cp.user_id FROM cast_gallery_images cgi
		JOIN cast_profiles cp ON cgi.cast_profile_id = cp.id
		WHERE cgi.id = $1
	`, imageID).Scan(&ownerID)

	if err == sql.ErrNoRows || ownerID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Delete image
	_, err = h.db.Exec("DELETE FROM cast_gallery_images WHERE id = $1", imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

func (h *CastHandler) GetCastBookings(c *gin.Context) {
	userID := c.GetInt("user_id")
	status := c.Query("status")
	
	query := `
		SELECT b.id, b.guest_id, b.cast_id, b.booking_date, b.start_time, 
		       b.duration_hours, b.location, b.amount, b.status,
		       b.created_at, u.name as guest_name, u.profile_image
		FROM bookings b
		JOIN users u ON b.guest_id = u.id
		WHERE b.cast_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND b.status = $2"
		args = append(args, status)
	}

	query += " ORDER BY b.booking_date DESC, b.start_time DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error getting cast bookings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	bookings := []gin.H{}
	for rows.Next() {
		var booking models.Booking
		var guestName string
		var profileImage sql.NullString

		err := rows.Scan(
			&booking.ID, &booking.GuestID, &booking.CastID,
			&booking.BookingDate, &booking.StartTime, &booking.DurationHours,
			&booking.Location, &booking.Amount, &booking.Status,
			&booking.CreatedAt, &guestName, &profileImage,
		)
		if err != nil {
			continue
		}

		bookingData := gin.H{
			"id":             booking.ID,
			"booking_date":   booking.BookingDate,
			"start_time":     booking.StartTime,
			"duration_hours": booking.DurationHours,
			"location":       booking.Location,
			"amount":         booking.Amount,
			"status":         booking.Status,
			"created_at":     booking.CreatedAt,
			"guest": gin.H{
				"id":            booking.GuestID,
				"name":          guestName,
				"profile_image": profileImage.String,
			},
		}
		bookings = append(bookings, bookingData)
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *CastHandler) RespondToBooking(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req models.BookingResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership and status
	var castID int
	var status models.BookingStatus
	var createdAt time.Time
	err = h.db.QueryRow(`
		SELECT cast_id, status, created_at FROM bookings WHERE id = $1
	`, bookingID).Scan(&castID, &status, &createdAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if castID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if status != models.BookingStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking already responded"})
		return
	}

	// Check if within 24 hours
	if time.Since(createdAt) > 24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Response time expired"})
		return
	}

	// Update booking status
	newStatus := models.BookingStatusAccepted
	var timestamp *time.Time
	now := time.Now()
	
	if !req.Accepted {
		newStatus = models.BookingStatusDeclined
		_, err = h.db.Exec(`
			UPDATE bookings SET status = $1, declined_at = $2, updated_at = $3
			WHERE id = $4
		`, newStatus, now, now, bookingID)
	} else {
		_, err = h.db.Exec(`
			UPDATE bookings SET status = $1, accepted_at = $2, updated_at = $3
			WHERE id = $4
		`, newStatus, now, now, bookingID)
		timestamp = &now
	}

	if err != nil {
		log.Printf("Error updating booking status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	// TODO: Send email notification to guest
	// TODO: If accepted, capture Stripe payment

	c.JSON(http.StatusOK, gin.H{
		"message":     "Booking response recorded",
		"status":      newStatus,
		"responded_at": timestamp,
	})
}

func (h *CastHandler) GetEarnings(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	// Get earnings summary
	var totalEarnings, thisMonthEarnings, pendingEarnings sql.NullFloat64
	
	err := h.db.QueryRow(`
		SELECT 
			SUM(CASE WHEN status = 'completed' THEN amount ELSE 0 END) as total_earnings,
			SUM(CASE WHEN status = 'completed' AND 
			    EXTRACT(YEAR FROM completed_at) = EXTRACT(YEAR FROM CURRENT_DATE) AND
			    EXTRACT(MONTH FROM completed_at) = EXTRACT(MONTH FROM CURRENT_DATE)
			    THEN amount ELSE 0 END) as this_month_earnings,
			SUM(CASE WHEN status = 'accepted' THEN amount ELSE 0 END) as pending_earnings
		FROM bookings
		WHERE cast_id = $1
	`, userID).Scan(&totalEarnings, &thisMonthEarnings, &pendingEarnings)

	if err != nil {
		log.Printf("Error getting earnings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get recent completed bookings
	rows, err := h.db.Query(`
		SELECT id, booking_date, amount, completed_at
		FROM bookings
		WHERE cast_id = $1 AND status = 'completed'
		ORDER BY completed_at DESC
		LIMIT 10
	`, userID)

	if err != nil {
		log.Printf("Error getting recent bookings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	recentBookings := []gin.H{}
	for rows.Next() {
		var id int
		var bookingDate, completedAt time.Time
		var amount float64
		
		if err := rows.Scan(&id, &bookingDate, &amount, &completedAt); err == nil {
			recentBookings = append(recentBookings, gin.H{
				"id":           id,
				"booking_date": bookingDate,
				"amount":       amount,
				"completed_at": completedAt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_earnings":      totalEarnings.Float64,
		"this_month_earnings": thisMonthEarnings.Float64,
		"pending_earnings":    pendingEarnings.Float64,
		"recent_bookings":     recentBookings,
	})
}