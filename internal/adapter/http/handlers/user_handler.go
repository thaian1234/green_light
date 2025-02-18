package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/thaian1234/green_light/internal/core/domain"
	"github.com/thaian1234/green_light/internal/core/ports"
	"github.com/thaian1234/green_light/pkg/logger"
	"github.com/thaian1234/green_light/pkg/util"
)

type UserHandler struct {
	wg            *sync.WaitGroup
	userService   ports.UserService
	mailerService ports.MailerService
}

func NewUserHandler(wg *sync.WaitGroup, userService ports.UserService, mailService ports.MailerService) *UserHandler {
	return &UserHandler{
		wg:            wg,
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
	util.Background(h.wg, func() {
		err = h.mailerService.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			logger.Error("failed to send welcome email", "msg", err)
		}
	})

	SendCreatedSuccess(ctx, Envelope{
		"user": user,
	})
}
