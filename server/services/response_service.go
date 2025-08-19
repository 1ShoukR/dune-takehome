package services

import (
	"context"
	"fmt"
	"time"

	"dune-takehome-server/database"
	"dune-takehome-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ResponseService struct {
	collection *mongo.Collection
}

func NewResponseService() *ResponseService {
	return &ResponseService{
		collection: database.Database.Collection("responses"),
	}
}

// CreateResponse saves a new form response
func (s *ResponseService) CreateResponse(formID primitive.ObjectID, req models.FormResponseRequest, ipAddress, userAgent string) (*models.FormUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := &models.FormUserResponse{
		ID:          primitive.NewObjectID(),
		FormID:      formID,
		Responses:   req.Responses,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		SubmittedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetFormResponses retrieves all responses for a specific form
func (s *ResponseService) GetFormResponses(formID primitive.ObjectID) ([]*models.FormUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Sort by submitted_at descending (most recent first)
	opts := options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, bson.M{"form_id": formID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var responses []*models.FormUserResponse
	if err = cursor.All(ctx, &responses); err != nil {
		return nil, err
	}

	return responses, nil
}

// GetResponseCount returns the total number of responses for a form
func (s *ResponseService) GetResponseCount(formID primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := s.collection.CountDocuments(ctx, bson.M{"form_id": formID})
	return count, err
}

// GetFormAnalytics generates analytics data for a form
func (s *ResponseService) GetFormAnalytics(form *models.Form) (*models.FormAnalytics, error) {
	responses, err := s.GetFormResponses(form.ID)
	if err != nil {
		return nil, err
	}

	analytics := &models.FormAnalytics{
		FormID:         form.ID,
		FormTitle:      form.Title,
		TotalResponses: int64(len(responses)),
		FieldAnalytics: []models.FieldAnalytics{},
		CreatedAt:      time.Now(),
	}

	// Generate analytics for each field
	for _, field := range form.Fields {
		fieldAnalytics := s.generateFieldAnalytics(field, responses)
		analytics.FieldAnalytics = append(analytics.FieldAnalytics, fieldAnalytics)
	}

	return analytics, nil
}

// generateFieldAnalytics creates analytics data for a specific field
func (s *ResponseService) generateFieldAnalytics(field models.FormField, responses []*models.FormUserResponse) models.FieldAnalytics {
	analytics := models.FieldAnalytics{
		FieldID:       field.ID,
		FieldLabel:    field.Label,
		FieldType:     string(field.Type),
		ResponseCount: 0,
		Data:          make(map[string]interface{}),
	}

	switch field.Type {
	case models.FieldTypeText, models.FieldTypeTextarea, models.FieldTypeEmail:
		analytics.Data = s.analyzeTextField(field.ID, responses)
	case models.FieldTypeNumber:
		analytics.Data = s.analyzeNumberField(field.ID, responses)
	case models.FieldTypeSelect, models.FieldTypeRadio:
		analytics.Data = s.analyzeChoiceField(field.ID, responses)
	case models.FieldTypeCheckbox:
		analytics.Data = s.analyzeCheckboxField(field.ID, responses)
	case models.FieldTypeRating:
		analytics.Data = s.analyzeRatingField(field.ID, responses)
	}

	// Count non-empty responses
	for _, response := range responses {
		if value, exists := response.Responses[field.ID]; exists && value != nil && value != "" {
			analytics.ResponseCount++
		}
	}

	return analytics
}

// analyzeTextField analyzes text-based fields
func (s *ResponseService) analyzeTextField(fieldID string, responses []*models.FormUserResponse) map[string]interface{} {
	data := make(map[string]interface{})
	var totalLength int
	var responseCount int

	for _, response := range responses {
		if value, exists := response.Responses[fieldID]; exists {
			if str, ok := value.(string); ok && str != "" {
				totalLength += len(str)
				responseCount++
			}
		}
	}

	if responseCount > 0 {
		data["average_length"] = float64(totalLength) / float64(responseCount)
	} else {
		data["average_length"] = 0
	}
	data["response_count"] = responseCount

	return data
}

// analyzeNumberField analyzes numeric fields
func (s *ResponseService) analyzeNumberField(fieldID string, responses []*models.FormUserResponse) map[string]interface{} {
	data := make(map[string]interface{})
	var total float64
	var count int
	var min, max float64
	var hasValues bool

	for _, response := range responses {
		if value, exists := response.Responses[fieldID]; exists {
			var num float64
			switch v := value.(type) {
			case float64:
				num = v
			case int:
				num = float64(v)
			default:
				continue
			}

			if !hasValues {
				min = num
				max = num
				hasValues = true
			} else {
				if num < min {
					min = num
				}
				if num > max {
					max = num
				}
			}
			total += num
			count++
		}
	}

	if count > 0 {
		data["average"] = total / float64(count)
		data["min"] = min
		data["max"] = max
	} else {
		data["average"] = 0
		data["min"] = 0
		data["max"] = 0
	}
	data["response_count"] = count

	return data
}

// analyzeChoiceField analyzes single-choice fields (select, radio)
func (s *ResponseService) analyzeChoiceField(fieldID string, responses []*models.FormUserResponse) map[string]interface{} {
	data := make(map[string]interface{})
	distribution := make(map[string]int)
	var totalResponses int

	for _, response := range responses {
		if value, exists := response.Responses[fieldID]; exists {
			if str, ok := value.(string); ok && str != "" {
				distribution[str]++
				totalResponses++
			}
		}
	}

	data["distribution"] = distribution
	data["response_count"] = totalResponses

	return data
}

// analyzeCheckboxField analyzes checkbox fields (multiple selections)
func (s *ResponseService) analyzeCheckboxField(fieldID string, responses []*models.FormUserResponse) map[string]interface{} {
	data := make(map[string]interface{})
	distribution := make(map[string]int)
	var totalResponses int

	for _, response := range responses {
		if value, exists := response.Responses[fieldID]; exists {
			if arr, ok := value.([]interface{}); ok {
				for _, item := range arr {
					if str, ok := item.(string); ok {
						distribution[str]++
					}
				}
				if len(arr) > 0 {
					totalResponses++
				}
			}
		}
	}

	data["distribution"] = distribution
	data["response_count"] = totalResponses

	return data
}

// analyzeRatingField analyzes rating fields
func (s *ResponseService) analyzeRatingField(fieldID string, responses []*models.FormUserResponse) map[string]interface{} {
	data := make(map[string]interface{})
	distribution := make(map[string]int)
	var total float64
	var count int

	for _, response := range responses {
		if value, exists := response.Responses[fieldID]; exists {
			var rating float64
			switch v := value.(type) {
			case float64:
				rating = v
			case int:
				rating = float64(v)
			default:
				continue
			}

			if rating >= 1 && rating <= 5 {
				ratingStr := fmt.Sprintf("%.0f", rating)
				distribution[ratingStr]++
				total += rating
				count++
			}
		}
	}

	if count > 0 {
		data["average_rating"] = total / float64(count)
	} else {
		data["average_rating"] = 0
	}
	data["distribution"] = distribution
	data["response_count"] = count

	return data
}
