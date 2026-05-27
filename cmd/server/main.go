package main

import (
	"habit-tracker/config"
	"habit-tracker/internal/app"
)

func main() {
	appConfig := config.NewConfig()

	appInstance := app.NewApp(appConfig)
	appInstance.Run()
}
