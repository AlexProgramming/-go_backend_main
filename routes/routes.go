package routes

import (
	"controllers"
	"github.com/gofiber/fiber/v2"
	"middleware"
)

func Router(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	// Post registers a route for POST methods that is used to submit an entity to the specified resource,
	//often causing a change in state or side effects on the server.

	//Use registers a middleware route that will match requests with
	//the provided prefix (which is optional and defaults to "/").
	//app.Use(func(c *fiber.Ctx) error {
	app.Use(middleware.VerifyAuthorization)

	// check for logged in before continue
	// authentication verification

	// update info
	// update password
	// log out
}
