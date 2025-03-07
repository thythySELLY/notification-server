package controllers

import (
	"net/http"
	dto "notification-server/modules/user-delivery/dtos"
	"notification-server/modules/user-delivery/models"
	"notification-server/modules/user-delivery/services"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserDeliveryController struct {
	service *services.UserDeliveryService
}

func NewUserDeliveryController(service *services.UserDeliveryService) *UserDeliveryController {
	return &UserDeliveryController{service: service}
}

func (c *UserDeliveryController) GetUserDeliveryList(ctx echo.Context) error {
	var query dto.GetUserDeliveryList

	if err := ctx.Bind(&query); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := c.service.GetUserDeliveryList(ctx.Request().Context(), query.Keyword, string(query.Status), query.Limit, query.PageToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserDeliveryController) CreateUserDelivery(ctx echo.Context) error {
	var req dto.CreateUserDelivery

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if strings.TrimSpace(req.Name) == "" || len(req.Name) > 100 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "name is required and must be under 100 characters"})
	}

	response, err := c.service.CreateUserDelivery(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *UserDeliveryController) UpdateUserDelivery(ctx echo.Context) error {
	id := ctx.Param("id")

	var req dto.UpdateUserDelivery

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	name := strings.TrimSpace(req.Name)

	id = strings.TrimSpace(id)

	if id == "" || name == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id and name are required"})
	}

	req.ID = id

	response, err := c.service.UpdateUserDeliveryService(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserDeliveryController) ChangeUserDeliveryStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.ChangeUserDeliveryStatus

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	req.ID = strings.TrimSpace(id)
	status := strings.TrimSpace(req.Status)

	if req.ID == "" || status == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id and status are required"})
	}

	if !models.IsValidStatus(status) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid status type"})
	}

	response, err := c.service.ChangeUserDeliveryStatus(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserDeliveryController) DeleteUserDelivery(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.DeleteUserDelivery
	req.ID = strings.TrimSpace(id)

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	response, err := c.service.DeleteUserDeliveryService(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}
