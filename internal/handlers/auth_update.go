package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		BirthDate string `json:"birth_date"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE users 
		SET name = $1, phone = $2, birth_date = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, req.Name, req.Phone, req.BirthDate, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *AuthHandler) UploadProfileImage(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	// In production, handle actual file upload to S3/Cloudinary
	imageURL := c.PostForm("image_url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL required"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE users 
		SET profile_image = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, imageURL, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile image updated successfully",
		"image_url": imageURL,
	})
}