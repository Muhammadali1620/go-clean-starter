package main

import (
	"flag"
	"log"

	_ "new_project/docs"
	"new_project/internal/app"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "http", "run mode: http | bot | hybrid | task-worker | scheduler")
}

// @title           NewProject API & Telegram Bot
// @version         1.0
// @description     Universal Backend API & Telegram Bot Server
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token in the format "Bearer <token>"

// @securityDefinitions.basic BasicAuth
func main() {
	flag.Parse()

	application := app.Bootstrap()
	defer application.Close()

	switch mode {
	case "http":
		app.RunHTTPServer(application)
	case "bot":
		app.RunBot(application)
	case "hybrid":
		app.RunHybrid(application)
	case "task-worker":
		app.RunWorker(application)
	case "scheduler":
		app.RunScheduler(application)
	default:
		log.Fatalf("❌ Unknown run mode: %s. Use: http | bot | hybrid | task-worker | scheduler", mode)
	}
}
