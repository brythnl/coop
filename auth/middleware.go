package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a Gin middleware for JWT authorization
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Missing authorization header"},
			)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid authorization header"},
			)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]
		token, err := jwt.ParseWithClaims(
			accessToken,
			&Claims{},
			func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid token: " + err.Error()},
			)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set(AuthorizationPayloadKey, claims)
		c.Next()
	}
}

// GetUserIDFromContext returns the user ID from the Gin context.
// To be used by handlers of protected routes.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return uuid.Nil, errors.New("authorization payload not found in context")
	}

	claims, ok := payload.(*Claims)
	if !ok {
		return uuid.Nil, errors.New("invalid payload type in context")
	}
	return claims.UserID, nil
}
