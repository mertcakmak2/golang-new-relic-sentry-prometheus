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

// CreateUser godoc
// @Summary Create User
// @Description Create User.
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User to be created"
// @Success 201 {object} domain.User "Returns created user"
// @Success 400 {object} domain.AppError "Returns error"
// @Router /api/v1/users [post]
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

// GetUserById godoc
// @Summary Get a user by ID
// @Description Retrieve a user using their ID from the database.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User "Returns user"
// @Success 404 {object} domain.AppError "Returns error"
// @Router /api/v1/users/{id} [get]
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

// UpdateUser godoc
// @Summary Update User
// @Description Update User.
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User to be updated"
// @Success 201 {object} domain.User "Returns updated user"
// @Success 400 {object} domain.AppError "Returns error"
// @Router /api/v1/users [put]
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

// DeleteUserById godoc
// @Summary Delete a user by ID
// @Description Delete a user using their ID from the database.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Success 500 {object} domain.AppError "Returns error"
// @Router /api/v1/users/{id} [delete]
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
