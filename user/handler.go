package user

import (
	"errors"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"go-app/domain"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	userUseCase domain.UserUseCase
	logger      *zap.Logger
}

func NewUserHandler(userUseCase domain.UserUseCase, logger *zap.Logger) *Handler {
	return &Handler{userUseCase: userUseCase, logger: logger}
}

func (h *Handler) CreateUser(c *gin.Context) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		var user domain.User

		if c.ShouldBind(&user) != nil {
			c.JSON(400, domain.NewBadRequestError("bad request"))
			return
		}

		createUser, err := h.userUseCase.CreateUser(user)
		if err != nil {
			hub.CaptureException(errors.New(err.Message))
			c.JSON(err.Code, err.AsMessageError())
			return
		}
		c.JSON(http.StatusCreated, createUser)
	}
}

func (h *Handler) GetUserById(c *gin.Context) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		idParam := c.Param("id")
		id, _ := strconv.ParseInt(idParam, 10, 64)

		user, err := h.userUseCase.GetUserById(uint(id))
		if err != nil {
			hub.CaptureException(errors.New(err.Message))
			c.JSON(err.Code, err.AsMessageError())
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) UpdateUser(c *gin.Context) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		var user domain.User
		if c.ShouldBind(&user) != nil {
			c.JSON(400, domain.NewBadRequestError("bad request"))
		}

		updatedUser, err := h.userUseCase.UpdateUser(user)
		if err != nil {
			hub.CaptureException(errors.New(err.Message))
			c.JSON(err.Code, err.AsMessageError())
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *Handler) DeleteUserById(c *gin.Context) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		idParam := c.Param("id")
		id, _ := strconv.ParseInt(idParam, 10, 64)

		err := h.userUseCase.DeleteUserById(uint(id))
		if err != nil {
			hub.CaptureException(errors.New(err.Message))
			c.JSON(err.Code, err.AsMessageError())
			return
		}
		c.Status(http.StatusNoContent)
	}
}
