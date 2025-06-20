package api

import (
	"context"
	"net/http"

	"github.com/brythnl/coop/auth"
	"github.com/brythnl/coop/db/sqlc"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store sqlc.Store, jwtSecret string) *gin.Engine {
	s := &Server{
		store:     store,
		jwtSecret: jwtSecret,
	}

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		if _, err := s.store.HealthCheck(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "DB error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Healthy!"})
	})

	apiV1 := router.Group("/api/v1")

	// Authentication
	authRoutes := apiV1.Group("/auth")
	authRoutes.POST("/register", s.RegisterUser)
	authRoutes.POST("/login", s.LoginUser)

	// Protected
	protected := apiV1.Group("").Use(auth.AuthMiddleware(s.jwtSecret))
	protected.GET("/me", s.getMeHandler)

	return router
}
