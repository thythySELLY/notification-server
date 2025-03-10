package controllers

import (
	"net/http"
	dto "notification-server/modules/connection/dtos"
	"notification-server/modules/connection/models"
	"notification-server/modules/connection/services"
	"strings"

	"github.com/labstack/echo/v4"
)

type ConnectionController struct {
	service *services.ConnectionService
}

func NewConnectionController(service *services.ConnectionService) *ConnectionController {
	return &ConnectionController{
		service: service,
	}
}

func (c *ConnectionController) CreateConnection(ctx echo.Context) error {
	var query dto.CreateConnection

	if err := ctx.Bind(&query); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if query.UserDeliveryServerId == "" || query.UserDeliveryServerWebHookUrl == "" || query.WebviewServerId == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "All fields must be non-empty"})
	}

	var connectionDto dto.CreateConnection = dto.CreateConnection{
		UserDeliveryServerId:         query.UserDeliveryServerId,
		UserDeliveryServerWebHookUrl: query.UserDeliveryServerWebHookUrl,
		WebviewServerId:              query.WebviewServerId,
	}

	response, err := c.service.CreateConnection(ctx.Request().Context(), connectionDto)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *ConnectionController) GetConnections(ctx echo.Context) error {
	var query dto.GetConnections

	if err := ctx.Bind(&query); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if query.Status != "active" && query.Status != "inactive" {
		query.Status = "inactive"
	}

	response, err := c.service.GetConnections(ctx.Request().Context(), query)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ConnectionController) UpdateWebHookUrl(ctx echo.Context) error {
	id := ctx.Param("id")

	var dto dto.UpdateUserDelivery
	dto.ID = id
	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "ID cannot be empty"})
	}
	if dto.UserDeliveryServerWebHookUrl == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "userDeliveryWebHookUrl cannot be empty"})
	}

	err := c.service.UpdateWebHookUrl(ctx.Request().Context(), dto)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Webhook URL updated successfully"})
}

func (c *ConnectionController) ChangeConnectionStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.ChangeConnectionStatus

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

	response, err := c.service.ChangeConnectionStatus(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *ConnectionController) DeleteConnection(ctx echo.Context) error {
	id := ctx.Param("id")

	var req dto.DeleteConnection
	req.ID = strings.TrimSpace(id)

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "ID cannot be empty"})
	}

	err := c.service.DeleteConnection(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Connection deleted successfully"})
}
