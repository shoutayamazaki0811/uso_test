package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/uso/uso/internal/database"
    "github.com/uso/uso/internal/models"
)

// ReportHandler handles report/block functionality
type ReportHandler struct {
    DB *database.DB
}

// ReportUser handles user reporting
func (h *ReportHandler) ReportUser(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    var req struct {
        ReportedUserID int64  `json:"reported_user_id" binding:"required"`
        Reason         string `json:"reason" binding:"required"`
        Description    string `json:"description"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    // Create report
    report := &models.Report{
        ReporterID:     userID,
        ReportedUserID: req.ReportedUserID,
        Reason:         req.Reason,
        Description:    req.Description,
        Status:         "pending",
    }
    
    if err := h.DB.CreateReport(report); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Report submitted successfully",
        "report_id": report.ID,
    })
}

// BlockUser handles user blocking
func (h *ReportHandler) BlockUser(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    var req struct {
        BlockedUserID int64 `json:"blocked_user_id" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    // Create block record
    block := &models.Block{
        UserID:        userID,
        BlockedUserID: req.BlockedUserID,
    }
    
    if err := h.DB.CreateBlock(block); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "User blocked successfully",
    })
}

// UnblockUser handles user unblocking
func (h *ReportHandler) UnblockUser(c *gin.Context) {
    userID := c.GetInt64("user_id")
    blockedUserID := c.Param("id")
    
    if err := h.DB.RemoveBlock(userID, blockedUserID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "User unblocked successfully",
    })
}

// GetBlockedUsers returns list of blocked users
func (h *ReportHandler) GetBlockedUsers(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    blockedUsers, err := h.DB.GetBlockedUsers(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blocked users"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "blocked_users": blockedUsers,
    })
}

// GetReports returns user's reports (admin only)
func (h *ReportHandler) GetReports(c *gin.Context) {
    // Check if user is admin
    userType := c.GetString("user_type")
    if userType != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        return
    }
    
    status := c.Query("status") // pending, reviewed, resolved
    
    reports, err := h.DB.GetReports(status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "reports": reports,
    })
}

// UpdateReportStatus updates report status (admin only)
func (h *ReportHandler) UpdateReportStatus(c *gin.Context) {
    // Check if user is admin
    userType := c.GetString("user_type")
    if userType != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        return
    }
    
    reportID := c.Param("id")
    
    var req struct {
        Status     string `json:"status" binding:"required"` // reviewed, resolved
        AdminNotes string `json:"admin_notes"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    if err := h.DB.UpdateReportStatus(reportID, req.Status, req.AdminNotes); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Report status updated successfully",
    })
}