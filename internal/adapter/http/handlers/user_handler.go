package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
)

type UserHandler struct {
	userService   ports.UserService
	mailerService ports.MailerService
}

func NewUserHandler(userService ports.UserService, mailService ports.MailerService) *UserHandler {
	return &UserHandler{
		userService:   userService,
		mailerService: mailService,
	}
}

type (
	registerUserRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
)

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Activated: false,
	}
	err := user.Password.Set(req.Password)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	if err = user.ValidateUser(); err != nil {
		HandleValidationError(ctx, err)
		return
	}
	if err = h.userService.CreateUser(ctx, user); err != nil {
		HandleError(ctx, err)
		return
	}
	err = h.mailerService.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	SendCreatedSuccess(ctx, Envelope{
		"user": user,
	})
}
