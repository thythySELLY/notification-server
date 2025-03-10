package api

import (
	"net/http"
	"notification-server/config"
	"notification-server/middlewares"
	connectionControllers "notification-server/modules/connection/controllers"
	connectionRepositories "notification-server/modules/connection/repositories"
	connectionServices "notification-server/modules/connection/services"
	userDeliveryControllers "notification-server/modules/user-delivery/controllers"
	userDeliveryRepositories "notification-server/modules/user-delivery/repositories"
	userDeliveryServices "notification-server/modules/user-delivery/services"
	webviewControllers "notification-server/modules/webview-server/controllers"
	webviewRepositories "notification-server/modules/webview-server/repositories"
	webviewServices "notification-server/modules/webview-server/services"

	"context"
	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeRouter() *echo.Echo {
	e := echo.New()

	clientOptions := options.Client().ApplyURI(config.MongoDBConfig.URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	webviewRepo := webviewRepositories.NewWebviewRepository(client.Database(config.MongoDBConfig.Database), client)
	userDeliveryRepo := userDeliveryRepositories.NewUserDeliveryRepository(client.Database(config.MongoDBConfig.Database), client)
	connectionRepo := connectionRepositories.NewConnectionRepository(client.Database(config.MongoDBConfig.Database))

	webViewService := webviewServices.NewWebviewService(webviewRepo)
	userDeliveryService := userDeliveryServices.NewUserDeliveryService(userDeliveryRepo, connectionRepo, webviewRepo)
	connectionService := connectionServices.NewConnectionService(connectionRepo, userDeliveryRepo, webviewRepo)

	webViewController := webviewControllers.NewWebViewController(webViewService)
	userDeliveryController := userDeliveryControllers.NewUserDeliveryController(userDeliveryService)
	connectionController := connectionControllers.NewConnectionController(connectionService)

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

	e.POST("/connection/new", connectionController.CreateConnection)
	e.GET("/connections", connectionController.GetConnections)
	e.PATCH("/connection/:id/webhook", connectionController.UpdateWebHookUrl)
	e.PATCH("/connection/:id/status", connectionController.ChangeConnectionStatus)

	return e
}
