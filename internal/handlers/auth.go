package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uso/uso/config"
	"github.com/uso/uso/internal/database"
	"github.com/uso/uso/internal/models"
	"github.com/uso/uso/internal/utils"
)

type AuthHandler struct {
	db  *database.DB
	cfg *config.Config
}

func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is 18+
	if req.BirthDate != nil {
		age := time.Now().Year() - req.BirthDate.Year()
		if time.Now().YearDay() < req.BirthDate.YearDay() {
			age--
		}
		if age < 18 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Must be 18 or older to register"})
			return
		}
	}

	// Check if email already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Insert user
	var userID int
	err = tx.QueryRow(`
		INSERT INTO users (email, password_hash, user_type, name, phone, birth_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, req.Email, hashedPassword, req.UserType, req.Name, req.Phone, req.BirthDate).Scan(&userID)
	
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// If cast, create pending cast profile
	if req.UserType == models.UserTypeCast {
		_, err = tx.Exec(`
			INSERT INTO cast_profiles (user_id, hourly_rate, rank, service_areas, approval_status)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, 60.0, models.CastRankStandard, []string{}, models.ApprovalStatusPending)
		
		if err != nil {
			log.Printf("Error creating cast profile: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating cast profile"})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(userID, req.Email, string(req.UserType), h.cfg.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":        userID,
			"email":     req.Email,
			"name":      req.Name,
			"user_type": req.UserType,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	var user models.User
	err := h.db.QueryRow(`
		SELECT id, email, password_hash, user_type, name, phone, birth_date, profile_image
		FROM users WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.UserType,
		&user.Name, &user.Phone, &user.BirthDate, &user.ProfileImage,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else if err != nil {
		log.Printf("Error querying user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.UserType), h.cfg.JWTSecret)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"user_type":     user.UserType,
			"profile_image": user.ProfileImage,
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var profile models.UserProfile
	err := h.db.QueryRow(`
		SELECT u.id, u.email, u.user_type, u.name, u.phone, u.birth_date, u.profile_image,
		       COALESCE(AVG(r.rating), 0) as rating, COUNT(r.id) as review_count
		FROM users u
		LEFT JOIN reviews r ON r.reviewed_id = u.id
		WHERE u.id = $1
		GROUP BY u.id
	`, userID).Scan(
		&profile.ID, &profile.Email, &profile.UserType, &profile.Name,
		&profile.Phone, &profile.BirthDate, &profile.ProfileImage,
		&profile.Rating, &profile.ReviewCount,
	)

	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// If cast, get cast profile
	if profile.UserType == models.UserTypeCast {
		var castProfile models.CastProfile
		err = h.db.QueryRow(`
			SELECT id, user_id, bio, hourly_rate, rank, service_areas, approval_status, approved_at
			FROM cast_profiles WHERE user_id = $1
		`, userID).Scan(
			&castProfile.ID, &castProfile.UserID, &castProfile.Bio,
			&castProfile.HourlyRate, &castProfile.Rank, &castProfile.ServiceAreas,
			&castProfile.ApprovalStatus, &castProfile.ApprovedAt,
		)
		
		if err == nil {
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
		}
	}

	c.JSON(http.StatusOK, profile)
}