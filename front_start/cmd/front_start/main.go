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

	"front_start/internal/app/config"
	"front_start/internal/app/dsn"
	"front_start/internal/app/handler"
	"front_start/internal/app/redis"
	"front_start/internal/app/repository"
	"front_start/internal/pkg"
	"context"

	// Импортируем сгенерированную документацию
	_ "front_start/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Summary Health check
// @Description Проверка работоспособности API
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
}

func main() {
	router := gin.Default()
	
	// Добавляем Swagger UI маршрут ДО инициализации приложения
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Health check маршрут
	router.GET("/health", healthCheck)

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