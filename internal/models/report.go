package models

import (
    "time"
)

// Report represents a user report
type Report struct {
    ID             int64     `json:"id" db:"id"`
    ReporterID     int64     `json:"reporter_id" db:"reporter_id"`
    ReportedUserID int64     `json:"reported_user_id" db:"reported_user_id"`
    Reason         string    `json:"reason" db:"reason"`
    Description    string    `json:"description" db:"description"`
    Status         string    `json:"status" db:"status"` // pending, reviewed, resolved
    AdminNotes     string    `json:"admin_notes" db:"admin_notes"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
    
    // Relations
    Reporter     *User `json:"reporter,omitempty"`
    ReportedUser *User `json:"reported_user,omitempty"`
}

// Block represents a user block
type Block struct {
    ID            int64     `json:"id" db:"id"`
    UserID        int64     `json:"user_id" db:"user_id"`
    BlockedUserID int64     `json:"blocked_user_id" db:"blocked_user_id"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    
    // Relations
    BlockedUser *User `json:"blocked_user,omitempty"`
}