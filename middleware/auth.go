package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthentication is a middleware function for JWT authentication.
func JWTAuthentication(c *fiber.Ctx) error {
	// Retrieve the token from the request headers
	token := c.Get("x-auth-token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid token",
		})
	}

	// Parse the token and get the claims
	claims, err := jwtToken(token)
	if err != nil {
		return fmt.Errorf("unauthorize")
	}

	// Print claims for debugging purposes (optional)
	fmt.Println(claims)

	// Check for expiration of the token
	expiresStr := claims["expires"].(string)

	expires, err := time.Parse(time.RFC3339, expiresStr)
	if err != nil {
		return fmt.Errorf("invalid expiration time format")
	}

	if time.Now().After(expires) {
		return fmt.Errorf("token has expired")
	}

	// Continue to the next handler
	return c.Next()
}

// jwtToken parses the JWT token string and returns the claims.
func jwtToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		// Retrieve the secret from the environment variable
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET is not set")
		}

		// Return the secret key for validation
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("unauthorized: token is invalid")
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized: unable to parse claims")
	}

	return claims, nil
}
