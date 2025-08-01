package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uso/uso/internal/models"
)

func (h *BookingHandler) GetMessages(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Verify user is part of booking and booking is accepted
	var guestID, castID int
	var status models.BookingStatus
	err = h.db.QueryRow(`
		SELECT guest_id, cast_id, status FROM bookings WHERE id = $1
	`, bookingID).Scan(&guestID, &castID, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if userID != guestID && userID != castID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if status != models.BookingStatusAccepted && status != models.BookingStatusCompleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "Messages only available for accepted bookings"})
		return
	}

	// Get messages
	rows, err := h.db.Query(`
		SELECT m.id, m.sender_id, m.message, m.created_at, u.name
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.booking_id = $1
		ORDER BY m.created_at ASC
	`, bookingID)

	if err != nil {
		log.Printf("Error getting messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	messages := []gin.H{}
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.Message, &msg.CreatedAt, &msg.SenderName)
		if err != nil {
			continue
		}

		messages = append(messages, gin.H{
			"id":          msg.ID,
			"sender_id":   msg.SenderID,
			"sender_name": msg.SenderName,
			"message":     msg.Message,
			"created_at":  msg.CreatedAt,
			"is_mine":     msg.SenderID == userID,
		})
	}

	c.JSON(http.StatusOK, messages)
}

func (h *BookingHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req models.MessageCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user is part of booking and booking is accepted
	var guestID, castID int
	var status models.BookingStatus
	err = h.db.QueryRow(`
		SELECT guest_id, cast_id, status FROM bookings WHERE id = $1
	`, bookingID).Scan(&guestID, &castID, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if userID != guestID && userID != castID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if status != models.BookingStatusAccepted && status != models.BookingStatusCompleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "Messages only available for accepted bookings"})
		return
	}

	// Insert message
	var messageID int
	var createdAt string
	err = h.db.QueryRow(`
		INSERT INTO messages (booking_id, sender_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, bookingID, userID, req.Message).Scan(&messageID, &createdAt)

	if err != nil {
		log.Printf("Error sending message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Get sender name
	var senderName string
	h.db.QueryRow("SELECT name FROM users WHERE id = $1", userID).Scan(&senderName)

	c.JSON(http.StatusCreated, gin.H{
		"id":          messageID,
		"sender_id":   userID,
		"sender_name": senderName,
		"message":     req.Message,
		"created_at":  createdAt,
		"is_mine":     true,
	})
}

func (h *BookingHandler) CreateReview(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req models.ReviewCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify booking exists and is completed
	var guestID, castID int
	var status models.BookingStatus
	err := h.db.QueryRow(`
		SELECT guest_id, cast_id, status FROM bookings WHERE id = $1
	`, req.BookingID).Scan(&guestID, &castID, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if userID != guestID && userID != castID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if status != models.BookingStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only review completed bookings"})
		return
	}

	// Check if already reviewed
	var exists bool
	err = h.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM reviews WHERE booking_id = $1 AND reviewer_id = $2)
	`, req.BookingID, userID).Scan(&exists)

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Already reviewed this booking"})
		return
	}

	// Determine reviewed_id
	reviewedID := castID
	if userID == castID {
		reviewedID = guestID
	}

	// Insert review
	var reviewID int
	err = h.db.QueryRow(`
		INSERT INTO reviews (booking_id, reviewer_id, reviewed_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, req.BookingID, userID, reviewedID, req.Rating, req.Comment).Scan(&reviewID)

	if err != nil {
		log.Printf("Error creating review: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         reviewID,
		"booking_id": req.BookingID,
		"rating":     req.Rating,
		"comment":    req.Comment,
	})
}

func (h *BookingHandler) GetUserReviews(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	rows, err := h.db.Query(`
		SELECT r.id, r.rating, r.comment, r.created_at, u.name, u.profile_image
		FROM reviews r
		JOIN users u ON r.reviewer_id = u.id
		WHERE r.reviewed_id = $1
		ORDER BY r.created_at DESC
	`, userID)

	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	reviews := []gin.H{}
	for rows.Next() {
		var review models.Review
		var reviewerName string
		var profileImage sql.NullString

		err := rows.Scan(&review.ID, &review.Rating, &review.Comment, 
			&review.CreatedAt, &reviewerName, &profileImage)
		if err != nil {
			continue
		}

		reviews = append(reviews, gin.H{
			"id":              review.ID,
			"rating":          review.Rating,
			"comment":         review.Comment,
			"created_at":      review.CreatedAt,
			"reviewer_name":   reviewerName,
			"reviewer_image":  profileImage.String,
		})
	}

	// Get average rating
	var avgRating sql.NullFloat64
	var totalReviews int
	err = h.db.QueryRow(`
		SELECT AVG(rating), COUNT(*) FROM reviews WHERE reviewed_id = $1
	`, userID).Scan(&avgRating, &totalReviews)

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avgRating.Float64,
		"total_reviews":  totalReviews,
		"reviews":        reviews,
	})
}