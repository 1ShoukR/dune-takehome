package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FormUserResponse represents a user's response to a shared form
type FormUserResponse struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	FormID      primitive.ObjectID     `json:"form_id" bson:"form_id"`
	Responses   map[string]interface{} `json:"responses" bson:"responses"`
	IPAddress   string                 `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	SubmittedAt time.Time              `json:"submitted_at" bson:"submitted_at"`
}

// FormResponseRequest represents the request payload for form submissions
type FormResponseRequest struct {
	Responses map[string]interface{} `json:"responses"`
}

// FormResponseSummary represents aggregated response data
type FormResponseSummary struct {
	FormID       primitive.ObjectID `json:"form_id"`
	TotalCount   int64              `json:"total_count"`
	LastResponse time.Time          `json:"last_response,omitempty"`
}

// FieldAnalytics represents analytics data for a specific field
type FieldAnalytics struct {
	FieldID       string                 `json:"field_id"`
	FieldLabel    string                 `json:"field_label"`
	FieldType     string                 `json:"field_type"`
	ResponseCount int64                  `json:"response_count"`
	Data          map[string]interface{} `json:"data"` // Flexible for different field types
}

// FormAnalytics represents complete analytics for a form
type FormAnalytics struct {
	FormID         primitive.ObjectID `json:"form_id"`
	FormTitle      string             `json:"form_title"`
	TotalResponses int64              `json:"total_responses"`
	FieldAnalytics []FieldAnalytics   `json:"field_analytics"`
	CreatedAt      time.Time          `json:"created_at"`
}
