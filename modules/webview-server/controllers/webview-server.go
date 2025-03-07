package controllers

import (
	"net/http"
	dto "notification-server/modules/webview-server/dtos"
	"notification-server/modules/webview-server/models"
	"notification-server/modules/webview-server/services"
	"strings"

	"github.com/labstack/echo/v4"
)

type WebViewController struct {
	service *services.WebViewService
}

func NewWebViewController(service *services.WebViewService) *WebViewController {
	return &WebViewController{service: service}
}

func (c *WebViewController) GetWebViewList(ctx echo.Context) error {
	var query dto.GetWebViewListQuery

	if err := ctx.Bind(&query); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response, err := c.service.GetWebviewListService(ctx.Request().Context(), query.Keyword, string(query.Status), query.Limit, query.PageToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *WebViewController) CreateWebView(ctx echo.Context) error {
	var req dto.CreateWebviewServer

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if strings.TrimSpace(req.Name) == "" || len(req.Name) > 100 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "name is required and must be under 100 characters"})
	}

	response, err := c.service.CreateWebviewService(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *WebViewController) UpdateWebView(ctx echo.Context) error {
	id := ctx.Param("id")

	var req dto.UpdateWebviewServer

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	name := strings.TrimSpace(req.Name)

	id = strings.TrimSpace(id)

	if id == "" || name == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id and name are required"})
	}

	req.ID = id

	response, err := c.service.UpdateWebviewService(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *WebViewController) ChangeWebViewStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.ChangeWebviewServerStatus

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

	response, err := c.service.ChangeWebviewStatus(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *WebViewController) DeleteWebview(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.DeleteWebviewServer
	req.ID = strings.TrimSpace(id)

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	response, err := c.service.DeleteWebviewService(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, response)
}
