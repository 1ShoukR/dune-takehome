package services

import (
	"context"
	"time"

	"dune-takehome-server/database"
	"dune-takehome-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormService struct {
	collection *mongo.Collection
}

func NewFormService() *FormService {
	return &FormService{
		collection: database.Database.Collection("forms"),
	}
}

// GetUserForms retrieves all forms for a specific user
func (s *FormService) GetUserForms(userID primitive.ObjectID, status *models.FormStatus) ([]*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	if status != nil {
		filter["status"] = *status
	}

	// Sort by updated_at descending (most recent first)
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var forms []*models.Form
	if err = cursor.All(ctx, &forms); err != nil {
		return nil, err
	}

	return forms, nil
}

// CreateForm creates a new form
func (s *FormService) CreateForm(userID primitive.ObjectID, req models.FormRequest) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = models.FormStatusDraft
	}

	form := &models.Form{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Fields:      req.Fields,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, form)
	if err != nil {
		return nil, err
	}

	return form, nil
}

// GetFormByID retrieves a form by ID
func (s *FormService) GetFormByID(formID primitive.ObjectID) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var form models.Form
	err := s.collection.FindOne(ctx, bson.M{"_id": formID}).Decode(&form)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Form not found
		}
		return nil, err
	}

	return &form, nil
}

// GetUserFormByID retrieves a form by ID that belongs to a specific user
func (s *FormService) GetUserFormByID(userID, formID primitive.ObjectID) (*models.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var form models.Form
	err := s.collection.FindOne(ctx, bson.M{
		"_id":     formID,
		"user_id": userID,
	}).Decode(&form)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Form not found
		}
		return nil, err
	}

	return &form, nil
}