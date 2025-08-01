package models

import (
	"database/sql/driver"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusAccepted  BookingStatus = "accepted"
	BookingStatusDeclined  BookingStatus = "declined"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

func (bs BookingStatus) Value() (driver.Value, error) {
	return string(bs), nil
}

func (bs *BookingStatus) Scan(value interface{}) error {
	*bs = BookingStatus(value.(string))
	return nil
}

type Booking struct {
	ID                   int           `json:"id"`
	GuestID              int           `json:"guest_id"`
	CastID               int           `json:"cast_id"`
	BookingDate          time.Time     `json:"booking_date"`
	StartTime            string        `json:"start_time"`
	DurationHours        int           `json:"duration_hours"`
	Location             string        `json:"location"`
	Amount               float64       `json:"amount"`
	Status               BookingStatus `json:"status"`
	StripePaymentIntentID *string      `json:"stripe_payment_intent_id,omitempty"`
	DeclinedAt           *time.Time    `json:"declined_at,omitempty"`
	AcceptedAt           *time.Time    `json:"accepted_at,omitempty"`
	CompletedAt          *time.Time    `json:"completed_at,omitempty"`
	CancelledAt          *time.Time    `json:"cancelled_at,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	Guest                *User         `json:"guest,omitempty"`
	Cast                 *User         `json:"cast,omitempty"`
}

type BookingCreate struct {
	CastID        int       `json:"cast_id" binding:"required"`
	BookingDate   time.Time `json:"booking_date" binding:"required"`
	StartTime     string    `json:"start_time" binding:"required"`
	DurationHours int       `json:"duration_hours" binding:"required,min=1"`
	Location      string    `json:"location" binding:"required"`
}

type BookingResponse struct {
	BookingID int    `json:"booking_id"`
	Accepted  bool   `json:"accepted"`
	Message   string `json:"message,omitempty"`
}

type Message struct {
	ID         int       `json:"id"`
	BookingID  int       `json:"booking_id"`
	SenderID   int       `json:"sender_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	SenderName string    `json:"sender_name,omitempty"`
}

type MessageCreate struct {
	BookingID int    `json:"booking_id" binding:"required"`
	Message   string `json:"message" binding:"required,max=1000"`
}

type Review struct {
	ID           int       `json:"id"`
	BookingID    int       `json:"booking_id"`
	ReviewerID   int       `json:"reviewer_id"`
	ReviewedID   int       `json:"reviewed_id"`
	Rating       int       `json:"rating"`
	Comment      *string   `json:"comment,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ReviewerName string    `json:"reviewer_name,omitempty"`
}

type ReviewCreate struct {
	BookingID int     `json:"booking_id" binding:"required"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Comment   *string `json:"comment" binding:"omitempty,max=500"`
}