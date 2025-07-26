package main

import (
	"fmt"
	"net/http"
	"os"

	"emtest/api-service/config"
	"emtest/api-service/db"
	handlers "emtest/api-service/handlers"
	"emtest/api-service/middleware"

	_ "emtest/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	cfg, err := config.LoadConfig("prod")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logrus.Error(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	err = db.InitDB(cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to init database connection: %s", err)
	}

	app.Use(cors.New())
	app.Use(middleware.Logger(logrus.StandardLogger()))

	app.Get("/swagger/*", swagger.HandlerDefault)

	v1 := app.Group("/api/v1")

	v1.Post("/subscriptions", handlers.CreateSubscription)
	v1.Get("/subscriptions", handlers.GetSubscriptions)
	v1.Put("/subscriptions", handlers.UpdateSubscription)
	v1.Delete("/subscriptions", handlers.DeleteSubscription)

	v1.Get("/subscriptions/calculate", handlers.CalculateTotalCost)

	logrus.Info("===============> Subscription CRUDL api <===============")
	logrus.Info("=> Project: " + "smth")
	logrus.Info("=> Host: " + cfg.Server.Host)
	logrus.Info(fmt.Sprintf("=> Port: %d", cfg.Server.Port))
	logrus.Info("=> Swagger: " + "/swagger")

	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

}
