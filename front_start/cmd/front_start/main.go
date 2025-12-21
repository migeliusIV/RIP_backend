// @title Quantum Tasks API
// @version 1.0
// @description API для управления квантовыми задачами и гейтами
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT токен в формате: "Bearer {token}"
package main

import (
	"fmt"
	"net/http"

	"context"
	"front_start/internal/app/config"
	"front_start/internal/app/dsn"
	"front_start/internal/app/handler"
	"front_start/internal/app/redis"
	"front_start/internal/app/repository"
	"front_start/internal/pkg"

	// Импортируем сгенерированную документацию
	_ "front_start/docs"

	// test lw6
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		// AllowOriginFunc: func(origin string) bool {
		// 	// 🔑 КРИТИЧЕСКИ ВАЖНО
		// 	if origin == "" || origin == "null" {
		// 		return true
		// 	}

		// 	switch origin {
		// 	case "http://localhost:5173":
		// 	case "http://127.0.0.1:5173":
		// 	case "https://migeliusiv.github.io":
		// 	case "":
		// 		return true
		// 	}

		// 	return false
		// },
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	// Добавляем Swagger UI маршрут ДО инициализации приложения
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Quantum Tasks API is running",
		})
	})

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	redisClient, errRedis := redis.New(context.Background(), conf.Redis)
	if errRedis != nil {
		logrus.Fatalf("error initializing redis: %v", errRedis)
	}

	hand := handler.NewHandler(rep, redisClient, &conf.JWT)

	application := pkg.NewApp(conf, router, hand)

	// Добавляем информационное сообщение
	fmt.Println("=== Quantum Tasks API ===")
	fmt.Println("Server started on: http://localhost:8080")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("=========================")

	application.RunApp()
}
