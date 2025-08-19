package middleware

import (
	"strings"

	"dune-takehome-server/services"
	"dune-takehome-server/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthRequired middleware validates JWT token and sets user info in context
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header must start with 'Bearer '",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is required",
			})
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Verify user still exists in database
		userService := services.NewUserService()
		user, err := userService.GetUserByID(claims.UserID)
		if err != nil || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Set user info in context
		c.Locals("userID", claims.UserID.Hex())
		c.Locals("userEmail", claims.Email)
		c.Locals("user", user)

		return c.Next()
	}
}

// GetCurrentUser helper function to get user from context
func GetCurrentUser(c *fiber.Ctx) *primitive.ObjectID {
	userID := c.Locals("userID")
	if userID == nil {
		return nil
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return nil
	}

	return &objectID
}
