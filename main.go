package main

import (
	"notification-server/api"

	"notification-server/config"
)

func main() {

	// Initialize MongoDB connection
	config.LoadEnv()
	config.InitMongoDB()
	config.InitRedis()

	e := api.InitializeRouter()
	e.Logger.Fatal(e.Start(":1323"))
}
