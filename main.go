package main

import (
	"db"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"routes"
)

type dbContainer struct {
	User    string
	Pass    string
	DbName  string
	Address string
}

func initEnvVars(cfg *dbContainer) {
	if err := godotenv.Load("environment_variables.env"); err != nil {
		log.Fatal("Failed to initialize")
	}

	cfg.User = os.Getenv("DBUSER")
	cfg.Pass = os.Getenv("DBPASS")
	cfg.DbName = os.Getenv("DBNAME")
	cfg.Address = os.Getenv("ADDRESS")
}

func main() {
	config := new(dbContainer)
	initEnvVars(config)

	db.Connect(fmt.Sprintf("%s:%s@/%s", config.User, config.Pass, config.DbName))

	app := fiber.New()

	// set up the apis
	routes.Router(app)

	panic(app.Listen(config.Address))
}
