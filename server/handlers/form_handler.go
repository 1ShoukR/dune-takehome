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
	formService     *services.FormService
	responseService *services.ResponseService
	wsService       *services.WebSocketService
}

func NewFormHandler(wsService *services.WebSocketService) *FormHandler {
	return &FormHandler{
		formService:     services.NewFormService(),
		responseService: services.NewResponseService(),
		wsService:       wsService,
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
	log.Printf("Looking for form with share_url: %s", c.Params("shareUrl"))

	if c.Params("shareUrl") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share URL is required",
		})
	}

	form, err := h.formService.GetFormByShareURL(c.Params("shareUrl"))
	if err != nil {
		log.Printf("Error fetching form: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve form",
		})
	}

	if form == nil {
		log.Printf("Form not found for share_url: %s", c.Params("shareUrl"))
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
	log.Printf("🔥 SubmitPublicFormResponse called with shareUrl: %s", c.Params("shareUrl"))

	shareURL := c.Params("shareUrl")
	if shareURL == "" {
		log.Printf("❌ Share URL is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Share URL is required",
		})
	}

	log.Printf("📝 Looking for form with shareUrl: %s", shareURL)

	form, err := h.formService.GetFormByShareURL(shareURL)
	if err != nil || form == nil || form.Status != models.FormStatusPublished {
		log.Printf("❌ Form not found or not published. Error: %v, Form: %v", err, form)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Form not found",
		})
	}

	log.Printf("✅ Form found: %s", form.Title)

	var req models.FormResponseRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("📋 Request body parsed, responses: %+v", req.Responses)

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	response, err := h.responseService.CreateResponse(form.ID, req, ipAddress, userAgent)
	if err != nil {
		log.Printf("❌ Failed to save response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save response",
		})
	}

	log.Printf("✅ Response saved successfully with ID: %s", response.ID.Hex())

	if h.wsService != nil {
		go func() {
			analytics, err := h.responseService.GetFormAnalytics(form)
			if err != nil {
				log.Printf("❌ Failed to generate analytics for broadcast: %v", err)
				return
			}
			h.wsService.BroadcastNewResponse(form.ID, analytics)
		}()
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Response submitted successfully",
		"response_id": response.ID.Hex(),
		"form_id":     form.ID.Hex(),
	})
}

// GetFormAnalytics returns analytics data for a form
func (h *FormHandler) GetFormAnalytics(c *fiber.Ctx) error {
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

	analytics, err := h.responseService.GetFormAnalytics(form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate analytics",
		})
	}

	return c.JSON(analytics)
}
