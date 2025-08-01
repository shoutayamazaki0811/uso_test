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

type AdminHandler struct {
	db  *database.DB
	cfg *config.Config
}

func NewAdminHandler(db *database.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{db: db, cfg: cfg}
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	// Get summary statistics
	var totalUsers, totalGuests, totalCasts, totalBookings int
	var totalRevenue sql.NullFloat64

	// Count users
	h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE user_type = 'guest'").Scan(&totalGuests)
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE user_type = 'cast'").Scan(&totalCasts)
	
	// Count bookings and revenue
	h.db.QueryRow("SELECT COUNT(*) FROM bookings").Scan(&totalBookings)
	h.db.QueryRow("SELECT SUM(amount) FROM bookings WHERE status = 'completed'").Scan(&totalRevenue)

	// Get recent activity
	recentBookings := []gin.H{}
	rows, err := h.db.Query(`
		SELECT b.id, b.booking_date, b.amount, b.status, 
		       g.name as guest_name, c.name as cast_name
		FROM bookings b
		JOIN users g ON b.guest_id = g.id
		JOIN users c ON b.cast_id = c.id
		ORDER BY b.created_at DESC
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var bookingDate time.Time
			var amount float64
			var status models.BookingStatus
			var guestName, castName string
			
			if err := rows.Scan(&id, &bookingDate, &amount, &status, &guestName, &castName); err == nil {
				recentBookings = append(recentBookings, gin.H{
					"id":           id,
					"booking_date": bookingDate,
					"amount":       amount,
					"status":       status,
					"guest_name":   guestName,
					"cast_name":    castName,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_users":    totalUsers,
			"total_guests":   totalGuests,
			"total_casts":    totalCasts,
			"total_bookings": totalBookings,
			"total_revenue":  totalRevenue.Float64,
		},
		"recent_bookings": recentBookings,
	})
}

func (h *AdminHandler) GetPendingCasts(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT u.id, u.name, u.email, u.created_at,
		       cp.id as profile_id, cp.bio, cp.rank, cp.service_areas
		FROM users u
		JOIN cast_profiles cp ON u.id = cp.user_id
		WHERE cp.approval_status = 'pending'
		ORDER BY u.created_at ASC
	`)
	if err != nil {
		log.Printf("Error getting pending casts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	pendingCasts := []gin.H{}
	for rows.Next() {
		var userID, profileID int
		var name, email string
		var createdAt time.Time
		var bio sql.NullString
		var rank models.CastRank
		var serviceAreas []string

		err := rows.Scan(&userID, &name, &email, &createdAt, 
			&profileID, &bio, &rank, &serviceAreas)
		if err != nil {
			continue
		}

		pendingCasts = append(pendingCasts, gin.H{
			"id":            userID,
			"name":          name,
			"email":         email,
			"created_at":    createdAt,
			"profile_id":    profileID,
			"bio":           bio.String,
			"rank":          rank,
			"service_areas": serviceAreas,
		})
	}

	c.JSON(http.StatusOK, pendingCasts)
}

func (h *AdminHandler) ApproveCast(c *gin.Context) {
	castID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cast ID"})
		return
	}

	// Update approval status
	result, err := h.db.Exec(`
		UPDATE cast_profiles 
		SET approval_status = $1, approved_at = $2, updated_at = $3
		WHERE user_id = $4 AND approval_status = 'pending'
	`, models.ApprovalStatusApproved, time.Now(), time.Now(), castID)

	if err != nil {
		log.Printf("Error approving cast: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve cast"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast not found or already processed"})
		return
	}

	// TODO: Send approval email

	c.JSON(http.StatusOK, gin.H{"message": "Cast approved successfully"})
}

func (h *AdminHandler) RejectCast(c *gin.Context) {
	castID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cast ID"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update approval status
	result, err := h.db.Exec(`
		UPDATE cast_profiles 
		SET approval_status = $1, updated_at = $2
		WHERE user_id = $3 AND approval_status = 'pending'
	`, models.ApprovalStatusRejected, time.Now(), castID)

	if err != nil {
		log.Printf("Error rejecting cast: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject cast"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cast not found or already processed"})
		return
	}

	// TODO: Send rejection email with reason

	c.JSON(http.StatusOK, gin.H{"message": "Cast rejected"})
}

func (h *AdminHandler) GetAllBookings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	query := `
		SELECT b.id, b.booking_date, b.start_time, b.duration_hours,
		       b.location, b.amount, b.status, b.created_at,
		       g.id, g.name, g.email,
		       c.id, c.name, c.email
		FROM bookings b
		JOIN users g ON b.guest_id = g.id
		JOIN users c ON b.cast_id = c.id
	`
	args := []interface{}{}
	
	if status != "" {
		query += " WHERE b.status = $1"
		args = append(args, status)
	}
	
	query += " ORDER BY b.created_at DESC"
	query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error getting bookings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	bookings := []gin.H{}
	for rows.Next() {
		var booking models.Booking
		var guestID, castID int
		var guestName, guestEmail, castName, castEmail string

		err := rows.Scan(
			&booking.ID, &booking.BookingDate, &booking.StartTime,
			&booking.DurationHours, &booking.Location, &booking.Amount,
			&booking.Status, &booking.CreatedAt,
			&guestID, &guestName, &guestEmail,
			&castID, &castName, &castEmail,
		)
		if err != nil {
			continue
		}

		bookings = append(bookings, gin.H{
			"id":             booking.ID,
			"booking_date":   booking.BookingDate,
			"start_time":     booking.StartTime,
			"duration_hours": booking.DurationHours,
			"location":       booking.Location,
			"amount":         booking.Amount,
			"status":         booking.Status,
			"created_at":     booking.CreatedAt,
			"guest": gin.H{
				"id":    guestID,
				"name":  guestName,
				"email": guestEmail,
			},
			"cast": gin.H{
				"id":    castID,
				"name":  castName,
				"email": castEmail,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"page":     page,
		"limit":    limit,
	})
}

func (h *AdminHandler) GetAnalytics(c *gin.Context) {
	// Get booking trends for last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	
	rows, err := h.db.Query(`
		SELECT DATE(booking_date) as date, COUNT(*) as count, SUM(amount) as revenue
		FROM bookings
		WHERE booking_date >= $1
		GROUP BY DATE(booking_date)
		ORDER BY date
	`, thirtyDaysAgo)

	bookingTrends := []gin.H{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date time.Time
			var count int
			var revenue sql.NullFloat64
			
			if err := rows.Scan(&date, &count, &revenue); err == nil {
				bookingTrends = append(bookingTrends, gin.H{
					"date":    date.Format("2006-01-02"),
					"count":   count,
					"revenue": revenue.Float64,
				})
			}
		}
	}

	// Get top performing casts
	topCasts := []gin.H{}
	castRows, err := h.db.Query(`
		SELECT c.id, c.name, COUNT(b.id) as booking_count, 
		       SUM(b.amount) as total_revenue, AVG(r.rating) as avg_rating
		FROM users c
		JOIN bookings b ON c.id = b.cast_id
		LEFT JOIN reviews r ON r.reviewed_id = c.id
		WHERE b.status = 'completed'
		GROUP BY c.id, c.name
		ORDER BY total_revenue DESC
		LIMIT 10
	`)
	
	if err == nil {
		defer castRows.Close()
		for castRows.Next() {
			var id int
			var name string
			var bookingCount int
			var totalRevenue, avgRating sql.NullFloat64
			
			if err := castRows.Scan(&id, &name, &bookingCount, &totalRevenue, &avgRating); err == nil {
				topCasts = append(topCasts, gin.H{
					"id":             id,
					"name":           name,
					"booking_count":  bookingCount,
					"total_revenue":  totalRevenue.Float64,
					"average_rating": avgRating.Float64,
				})
			}
		}
	}

	// Get popular service areas
	areaStats := []gin.H{}
	areaRows, err := h.db.Query(`
		SELECT sa.name, COUNT(DISTINCT cp.user_id) as cast_count
		FROM service_areas sa
		JOIN cast_profiles cp ON sa.name = ANY(cp.service_areas)
		WHERE cp.approval_status = 'approved'
		GROUP BY sa.name
		ORDER BY cast_count DESC
	`)
	
	if err == nil {
		defer areaRows.Close()
		for areaRows.Next() {
			var areaName string
			var castCount int
			
			if err := areaRows.Scan(&areaName, &castCount); err == nil {
				areaStats = append(areaStats, gin.H{
					"area":       areaName,
					"cast_count": castCount,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"booking_trends":     bookingTrends,
		"top_casts":         topCasts,
		"service_area_stats": areaStats,
	})
}