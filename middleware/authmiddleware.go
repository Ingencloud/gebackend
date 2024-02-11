package middleware

import (

	// "encoding/base64"

	"github.com/gofiber/fiber/v2"

	// "github.com/twilio/twilio-go/rest/api/v2010"

	// "strconv"
	"ge/utils"
)

// IsAuthenticated middleware
// IsAuthenticated middleware
// IsAuthenticated middleware

// IsAuthenticated middleware
func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	// Parse the JWT token
	userID, err := utils.ParseJwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Add the UserID to the request context
	c.Locals("UserID", userID)

	c.Set("Authorization", "Bearer "+cookie)

	return c.Next()
}

// 	// Optionally, you can pass the user ID to subsequent handlers if needed
// 	c.Locals("userId", userId)

// 	return c.Next()
// }

// func Isauthenticated(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")

// 	if _, err := utils.ParseJwt(cookie); err != nil {
// 		c.Status(fiber.StatusUnauthorized)
// 		return c.JSON(fiber.Map{
// 			"message": "unauthenticated",
// 		})
// 	}
// 	return c.Next()
// }
