package api

import (
	"context"
	"net/http"

	"github.com/brythnl/coop/db/sqlc"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store sqlc.Store, jwtSecret string) *gin.Engine {
	s := &Server{
		store:     store,
		jwtSecret: jwtSecret,
	}

	router := gin.Default()

	router.GET("/healthz", func(ctx *gin.Context) {
		if _, err := s.store.HealthCheck(context.Background()); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "db error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "healthy!"})
	})

	apiV1 := router.Group("/api/v1")

	authRoutes := apiV1.Group("/auth")
	authRoutes.POST("/register", s.RegisterUser)
	authRoutes.POST("/login", s.LoginUser)

	return router
}
