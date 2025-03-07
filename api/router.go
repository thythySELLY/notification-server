package api

import (
	"net/http"
	"notification-server/config"
	"notification-server/middlewares"
	webviewControllers "notification-server/modules/webview-server/controllers"
	webviewRepositories "notification-server/modules/webview-server/repositories"
	webviewServices "notification-server/modules/webview-server/services"
	userDeliveryControllers "notification-server/modules/user-delivery/controllers"
	userDeliveryRepositories "notification-server/modules/user-delivery/repositories"
	userDeliveryServices "notification-server/modules/user-delivery/services"

	"github.com/labstack/echo/v4"
)

func InitializeRouter() *echo.Echo {
	e := echo.New()

	webviewRepo := webviewRepositories.NewWebviewRepository(config.MongoDBClient.Database(config.MongoDBConfig.Database))
	webViewService := webviewServices.NewWebviewService(webviewRepo)
	webViewController := webviewControllers.NewWebViewController(webViewService)

	userDeliveryRepo := userDeliveryRepositories.NewUserDeliveryRepository(config.MongoDBClient.Database(config.MongoDBConfig.Database))
	userDeliveryService := userDeliveryServices.NewUserDeliveryService(userDeliveryRepo)
	userDeliveryController := userDeliveryControllers.NewUserDeliveryController(userDeliveryService)

	e.Use(middlewares.ValidateToken)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, This is Notification Server!")
	})

	e.GET("/webview-servers", webViewController.GetWebViewList)
	e.POST("/webview-server", webViewController.CreateWebView)
	e.PUT("/webview-server/:id", webViewController.UpdateWebView)
	e.PATCH("/webview-server/:id/status", webViewController.ChangeWebViewStatus)
	e.DELETE("/webview-server/:id", webViewController.DeleteWebview)

	e.GET("/user-deliveries", userDeliveryController.GetUserDeliveryList)
	e.POST("/user-delivery", userDeliveryController.CreateUserDelivery)
	e.PUT("/user-delivery/:id", userDeliveryController.UpdateUserDelivery)
	e.PATCH("/user-delivery/:id/status", userDeliveryController.ChangeUserDeliveryStatus)
	e.DELETE("/user-delivery/:id", userDeliveryController.DeleteUserDelivery)

	return e
}
