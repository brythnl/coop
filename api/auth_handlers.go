package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/brythnl/coop/auth"
	"github.com/brythnl/coop/db/sqlc"
	"github.com/brythnl/coop/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Server struct {
	store     sqlc.Store
	jwtSecret string
}

func (s *Server) RegisterUser(ctx *gin.Context) {
	var payload models.RegisterPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	params := sqlc.CreateUserParams{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}

	user, err := s.store.CreateUser(context.Background(), params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var payload models.LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user, err := s.store.GetUserByUsername(context.Background(), payload.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !auth.CheckPasswordHash(payload.Password, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
