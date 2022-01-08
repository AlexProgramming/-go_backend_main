package middleware

import (
	"github.com/gofiber/fiber/v2"
	"tokens"
)

func VerifyAuthorization(c *fiber.Ctx) error {
	cookie := c.Cookies("rest_cookie") // synchronize with authController - use property file

	if _, err := tokens.VerifyJWT(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)

		return c.JSON(fiber.Map{
			"response": "Unable to verify login",
		})
	}

	return c.Next()
}
