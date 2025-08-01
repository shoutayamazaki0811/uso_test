package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uso/uso/config"
	"github.com/uso/uso/internal/database"
	"github.com/uso/uso/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

type BookingHandler struct {
	db  *database.DB
	cfg *config.Config
}

func NewBookingHandler(db *database.DB, cfg *config.Config) *BookingHandler {
	return &BookingHandler{db: db, cfg: cfg}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req models.BookingCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify cast exists and is approved
	var castApprovalStatus models.ApprovalStatus
	var hourlyRate float64
	err := h.db.QueryRow(`
		SELECT cp.approval_status, cp.hourly_rate
		FROM users u
		JOIN cast_profiles cp ON u.id = cp.user_id
		WHERE u.id = $1 AND u.user_type = 'cast'
	`, req.CastID).Scan(&castApprovalStatus, &hourlyRate)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast not found"})
		return
	}

	if castApprovalStatus != models.ApprovalStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cast not approved"})
		return
	}

	// Check for booking conflicts
	var conflicts int
	err = h.db.QueryRow(`
		SELECT COUNT(*) FROM bookings
		WHERE cast_id = $1 
		AND booking_date = $2
		AND status IN ('pending', 'accepted')
		AND (
			($3::time >= start_time AND $3::time < start_time + (duration_hours || ' hours')::interval) OR
			(start_time >= $3::time AND start_time < $3::time + ($4 || ' hours')::interval)
		)
	`, req.CastID, req.BookingDate, req.StartTime, req.DurationHours).Scan(&conflicts)

	if conflicts > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cast already has a booking at this time"})
		return
	}

	// Calculate amount
	amount := hourlyRate * float64(req.DurationHours)

	// Create Stripe payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)), // Convert to cents
		Currency: stripe.String("usd"),
		Metadata: map[string]string{
			"guest_id": strconv.Itoa(userID),
			"cast_id":  strconv.Itoa(req.CastID),
		},
		CaptureMethod: stripe.String("manual"), // Don't capture until cast accepts
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment processing failed"})
		return
	}

	// Create booking
	var bookingID int
	err = h.db.QueryRow(`
		INSERT INTO bookings (guest_id, cast_id, booking_date, start_time, duration_hours, 
		                     location, amount, status, stripe_payment_intent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, userID, req.CastID, req.BookingDate, req.StartTime, req.DurationHours,
	   req.Location, amount, models.BookingStatusPending, pi.ID).Scan(&bookingID)

	if err != nil {
		log.Printf("Error creating booking: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// TODO: Send email notification to cast

	c.JSON(http.StatusCreated, gin.H{
		"booking_id":     bookingID,
		"amount":         amount,
		"status":         models.BookingStatusPending,
		"client_secret":  pi.ClientSecret,
		"payment_intent": pi.ID,
	})
}

func (h *BookingHandler) GetGuestBookings(c *gin.Context) {
	userID := c.GetInt("user_id")
	status := c.Query("status")
	
	query := `
		SELECT b.id, b.guest_id, b.cast_id, b.booking_date, b.start_time, 
		       b.duration_hours, b.location, b.amount, b.status,
		       b.created_at, u.name as cast_name, u.profile_image,
		       cp.rank, COALESCE(AVG(r.rating), 0) as rating
		FROM bookings b
		JOIN users u ON b.cast_id = u.id
		JOIN cast_profiles cp ON u.id = cp.user_id
		LEFT JOIN reviews r ON r.reviewed_id = u.id
		WHERE b.guest_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND b.status = $2"
		args = append(args, status)
	}

	query += " GROUP BY b.id, u.name, u.profile_image, cp.rank ORDER BY b.booking_date DESC, b.start_time DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error getting guest bookings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	bookings := []gin.H{}
	for rows.Next() {
		var booking models.Booking
		var castName string
		var profileImage sql.NullString
		var rank models.CastRank
		var rating float64

		err := rows.Scan(
			&booking.ID, &booking.GuestID, &booking.CastID,
			&booking.BookingDate, &booking.StartTime, &booking.DurationHours,
			&booking.Location, &booking.Amount, &booking.Status,
			&booking.CreatedAt, &castName, &profileImage, &rank, &rating,
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
			"cast": gin.H{
				"id":            booking.CastID,
				"name":          castName,
				"profile_image": profileImage.String,
				"rank":          rank,
				"rating":        rating,
			},
		}
		bookings = append(bookings, bookingData)
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	var guestName, castName string
	var guestImage, castImage sql.NullString

	err = h.db.QueryRow(`
		SELECT b.*, 
		       g.name as guest_name, g.profile_image as guest_image,
		       c.name as cast_name, c.profile_image as cast_image
		FROM bookings b
		JOIN users g ON b.guest_id = g.id
		JOIN users c ON b.cast_id = c.id
		WHERE b.id = $1 AND (b.guest_id = $2 OR b.cast_id = $2)
	`, bookingID, userID).Scan(
		&booking.ID, &booking.GuestID, &booking.CastID,
		&booking.BookingDate, &booking.StartTime, &booking.DurationHours,
		&booking.Location, &booking.Amount, &booking.Status,
		&booking.StripePaymentIntentID, &booking.DeclinedAt, &booking.AcceptedAt,
		&booking.CompletedAt, &booking.CancelledAt, &booking.CreatedAt,
		&booking.UpdatedAt, &guestName, &guestImage, &castName, &castImage,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
			"profile_image": guestImage.String,
		},
		"cast": gin.H{
			"id":            booking.CastID,
			"name":          castName,
			"profile_image": castImage.String,
		},
	})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Verify ownership and status
	var guestID int
	var status models.BookingStatus
	var paymentIntentID sql.NullString
	err = h.db.QueryRow(`
		SELECT guest_id, status, stripe_payment_intent_id 
		FROM bookings WHERE id = $1
	`, bookingID).Scan(&guestID, &status, &paymentIntentID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if guestID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if status != models.BookingStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only cancel pending bookings"})
		return
	}

	// Cancel Stripe payment intent
	if paymentIntentID.Valid {
		_, err = paymentintent.Cancel(paymentIntentID.String, nil)
		if err != nil {
			log.Printf("Error cancelling payment intent: %v", err)
		}
	}

	// Update booking status
	_, err = h.db.Exec(`
		UPDATE bookings 
		SET status = $1, cancelled_at = $2, updated_at = $3
		WHERE id = $4
	`, models.BookingStatusCancelled, time.Now(), time.Now(), bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

func (h *BookingHandler) CompleteBooking(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Verify booking exists and is accepted
	var castID, guestID int
	var status models.BookingStatus
	var bookingDate time.Time
	var startTime string
	var durationHours int

	err = h.db.QueryRow(`
		SELECT cast_id, guest_id, status, booking_date, start_time, duration_hours
		FROM bookings WHERE id = $1
	`, bookingID).Scan(&castID, &guestID, &status, &bookingDate, &startTime, &durationHours)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Only cast can mark as completed
	if castID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only cast can complete booking"})
		return
	}

	if status != models.BookingStatusAccepted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking must be accepted first"})
		return
	}

	// Check if booking time has passed
	bookingEndTime := bookingDate.Add(time.Duration(durationHours) * time.Hour)
	if time.Now().Before(bookingEndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot complete booking before end time"})
		return
	}

	// Update booking status
	_, err = h.db.Exec(`
		UPDATE bookings 
		SET status = $1, completed_at = $2, updated_at = $3
		WHERE id = $4
	`, models.BookingStatusCompleted, time.Now(), time.Now(), bookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete booking"})
		return
	}

	// TODO: Send review request emails

	c.JSON(http.StatusOK, gin.H{"message": "Booking completed successfully"})
}

func (h *BookingHandler) HandleStripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), h.cfg.StripeWebhookSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
			return
		}
		
		// Update booking payment status
		log.Printf("Payment succeeded for intent: %s", paymentIntent.ID)
		
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
			return
		}
		
		// Handle failed payment
		log.Printf("Payment failed for intent: %s", paymentIntent.ID)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}