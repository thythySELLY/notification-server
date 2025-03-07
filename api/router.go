package api

import (
	"net/http"
	"notification-server/config"
	"notification-server/middlewares"
	"notification-server/modules/webview-server/controllers"
	"notification-server/modules/webview-server/repositories"
	"notification-server/modules/webview-server/services"

	"github.com/labstack/echo/v4"
)

func InitializeRouter() *echo.Echo {
	e := echo.New()

	webviewRepo := repositories.NewWebviewRepository(config.MongoDBClient.Database(config.MongoDBConfig.Database))
	webViewService := services.NewWebviewService(webviewRepo)
	webViewController := controllers.NewWebViewController(webViewService)

	e.Use(middlewares.ValidateToken)

	// Health Check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, This is Notification Server!")
	})

	e.GET("/webview-servers", webViewController.GetWebViewList, middlewares.ValidateToken)
	e.POST("/webview-server", webViewController.CreateWebView, middlewares.ValidateToken)
	e.PUT("/webview-server/:id", webViewController.UpdateWebView, middlewares.ValidateToken)
	e.PATCH("/webview-server/:id/status", webViewController.ChangeWebViewStatus, middlewares.ValidateToken)
	e.DELETE("/webview-server/:id", webViewController.DeleteWebview, middlewares.ValidateToken)

	return e
}
