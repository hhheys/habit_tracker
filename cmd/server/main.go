package main

import (
	"habit-tracker/config"
	"habit-tracker/internal/app"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	appConfig := config.NewConfig()

	appInstance := app.NewApp(appConfig)
	appInstance.Run()
}
