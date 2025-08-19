package handlers

import (
	"log"
	"strings"

	"dune-takehome-server/models"
	"dune-takehome-server/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormHandler struct {
	formService *services.FormService
}

func NewFormHandler() *FormHandler {
	return &FormHandler{
		formService: services.NewFormService(),
	}
}

// GetUserForms retrieves all forms for the authenticated user
func (h *FormHandler) GetUserForms(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check for status filter
	var statusFilter *models.FormStatus
	if statusQuery := c.Query("status"); statusQuery != "" {
		status := models.FormStatus(statusQuery)
		statusFilter = &status
	}

	// Get forms
	forms, err := h.formService.GetUserForms(userID, statusFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve forms",
		})
	}

	// Convert to response format
	var formResponses []models.FormResponse
	for _, form := range forms {
		formResponses = append(formResponses, form.ToResponse())
	}

	return c.JSON(fiber.Map{
		"forms": formResponses,
		"count": len(formResponses),
	})
}

// CreateForm creates a new form for the authenticated user
func (h *FormHandler) CreateForm(c *fiber.Ctx) error {
	// Get user ID from auth middleware
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Parse request body
	var req models.FormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if strings.TrimSpace(req.Title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Form title is required",
		})
	}

	// Create form
	form, err := h.formService.CreateForm(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create form",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(form.ToResponse())
}

// GetFormByID retrieves a specific form
func (h *FormHandler) GetFormByID(c *fiber.Ctx) error {
	// Get form ID from URL params
	formIDStr := c.Params("id")
	formID, err := primitive.ObjectIDFromHex(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form ID",
		})
	}

	// Get user ID from auth middleware
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get form (ensure it belongs to the user)
	form, err := h.formService.GetUserFormByID(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve form",
		})
	}

	if form == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	return c.JSON(form.ToResponse())
}

// UpdateForm updates an existing form
func (h *FormHandler) UpdateForm(c *fiber.Ctx) error {
	formIDStr := c.Params("id")
	formID, err := primitive.ObjectIDFromHex(formIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form ID",
		})
	}

	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.FormRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if strings.TrimSpace(req.Title) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Form title is required",
		})
	}

	form, err := h.formService.UpdateForm(userID, formID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update form",
		})
	}

	if form == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	return c.JSON(form.ToResponse())
}

// GetPublicForm retrieves a form by share URL (no auth required)
func (h *FormHandler) GetPublicForm(c *fiber.Ctx) error {
	shareURL := c.Params("shareUrl")
	log.Printf("Looking for form with share_url: %s", shareURL)

	if shareURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share URL is required",
		})
	}

	form, err := h.formService.GetFormByShareURL(shareURL)
	if err != nil {
		log.Printf("Error fetching form: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve form",
		})
	}

	if form == nil {
		log.Printf("Form not found for share_url: %s", shareURL)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	if form.Status != models.FormStatusPublished {
		log.Printf("Form found but not published. Status: %s", form.Status)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	log.Printf("Form found successfully: %s", form.Title)
	return c.JSON(form.ToResponse())
}

// SubmitPublicFormResponse handles form submissions (no auth required)
func (h *FormHandler) SubmitPublicFormResponse(c *fiber.Ctx) error {
	shareURL := c.Params("shareUrl")
	if shareURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share URL is required",
		})
	}

	form, err := h.formService.GetFormByShareURL(shareURL)
	if err != nil || form == nil || form.Status != models.FormStatusPublished {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	var req struct {
		Responses map[string]interface{} `json:"responses"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// TODO: Save response to database
	// For now, just return success
	return c.JSON(fiber.Map{
		"message": "Response submitted successfully",
		"form_id": form.ID.Hex(),
	})
}
