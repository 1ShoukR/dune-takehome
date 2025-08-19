package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FormStatus represents the status of a form
type FormStatus string

const (
	FormStatusDraft     FormStatus = "draft"
	FormStatusPublished FormStatus = "published"
)

// FieldType represents the type of form field
type FieldType string

const (
	FieldTypeText       FieldType = "text"
	FieldTypeTextarea   FieldType = "textarea"
	FieldTypeEmail      FieldType = "email"
	FieldTypeNumber     FieldType = "number"
	FieldTypeSelect     FieldType = "select"
	FieldTypeRadio      FieldType = "radio"
	FieldTypeCheckbox   FieldType = "checkbox"
	FieldTypeRating     FieldType = "rating"
)

// FormField represents a field in a form
type FormField struct {
	ID          string            `json:"id" bson:"id"`
	Type        FieldType         `json:"type" bson:"type"`
	Label       string            `json:"label" bson:"label"`
	Placeholder string            `json:"placeholder,omitempty" bson:"placeholder,omitempty"`
	Required    bool              `json:"required" bson:"required"`
	Options     []string          `json:"options,omitempty" bson:"options,omitempty"` // For select, radio, checkbox
	Validation  map[string]string `json:"validation,omitempty" bson:"validation,omitempty"`
	Order       int               `json:"order" bson:"order"`
}

// Form represents a form document
type Form struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Fields      []FormField        `json:"fields" bson:"fields"`
	Status      FormStatus         `json:"status" bson:"status"`
	ShareURL    string             `json:"share_url,omitempty" bson:"share_url,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// FormRequest represents the request payload for creating/updating forms
type FormRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Fields      []FormField `json:"fields"`
	Status      FormStatus  `json:"status,omitempty"`
}

// FormResponse represents the response payload for form data
type FormResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description,omitempty"`
	Fields      []FormField        `json:"fields"`
	Status      FormStatus         `json:"status"`
	ShareURL    string             `json:"share_url,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToResponse converts Form to FormResponse
func (f *Form) ToResponse() FormResponse {
	return FormResponse{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
		Fields:      f.Fields,
		Status:      f.Status,
		ShareURL:    f.ShareURL,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}