package handler

import (
	"front_start/internal/app/config"
	appredis "front_start/internal/app/redis"
	"front_start/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler структура обработчиков API
// @Description Основная структура содержащая зависимости обработчиков
type Handler struct {
	Repository *repository.Repository
	Redis      *appredis.Client
	JWTConfig  *config.JWTConfig
}

// NewHandler создает новый экземпляр Handler
// @Description Конструктор для создания экземпляра Handler с зависимостями
func NewHandler(r *repository.Repository, redis *appredis.Client, jwtConfig *config.JWTConfig) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redis,
		JWTConfig:  jwtConfig,
	}
}

// RegisterHandler регистрирует все маршруты API
// @Description Функция для регистрации всех маршрутов приложения с группировкой по правам доступа
func (handler *Handler) RegisterHandler(r *gin.Engine) {
	// Публичные эндпоинты (без авторизации)
	r.POST("/login", handler.Login)
	r.POST("/users", handler.Register)
	// gates
	r.GET("/api/gates", handler.ApiGatesList)
	r.GET("/api/gates/:id", handler.ApiGetGateByID)
	// HTML
	r.GET("/IBM", handler.GetGates)
	r.GET("/gate_property/:id", handler.GetGateByID)

	//internal := r.Group("/internal")
	//{
	//internal.PUT("/quantum_task/res", handler.SetQuantumTaskResult)
	//}

	// Эндпоинты, доступные только модераторам
	moderator := r.Group("/")
	moderator.Use(handler.AuthMiddleware, handler.ModeratorMiddleware)
	{
		// Gates (создание, изменение, удаление)
		moderator.POST("/api/gates", handler.ApiAddGate)
		moderator.PUT("/api/gates/:id", handler.ApiUpdateGate)
		moderator.DELETE("/api/gates/:id", handler.ApiDeleteGate)
		moderator.POST("/api/gates/:id/image", handler.ApiUploadGatesImage)

		// QuantumTasks (завершение/отклонение)
		moderator.PUT("/api/quantum_tasks/:id/resolve", handler.ApiResolveQTask)
	}

	// Эндпоинты, доступные всем авторизованным пользователям
	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware)
	{
		// API Users
		auth.POST("/api/auth/logout", handler.Logout)
		auth.GET("/api/users/me", handler.ApiMe)
		auth.PUT("/api/users/me", handler.ApiUpdateMe)

		// Связи задачи-гейты (many-to-many)
		auth.DELETE("/api/tasks/:task_id/services/:service_id", handler.ApiRemoveGateFromTask)
		auth.PUT("/api/tasks/:task_id/services/:service_id", handler.ApiUpdateDegrees)

		// HTML
		auth.GET("/quantum_task/:id", handler.GetTask)
		auth.POST("/quantum_task/add/gate/:id_gate", handler.AddGateToTask)
		auth.POST("/quantum_task/:task_id/delete", handler.DeleteTask)

		// API Gates
		auth.POST("/api/draft/gates/:id", handler.ApiAddGateToDraft)

		// API Quantum tasks
		auth.GET("/api/quantum_task/current", handler.ApiGetCurrQTask)
		auth.GET("/api/quantum_tasks", handler.ApiListQTasks)
		auth.GET("/api/quantum_tasks/:id", handler.ApiGetQTaskByID)
		auth.PUT("/api/quantum_tasks/:id", handler.ApiUpdateQTask)
		auth.PUT("/api/quantum_tasks/:id/form", handler.ApiFormQTask)
		auth.DELETE("/api/quantum_tasks/:id", handler.ApiDeleteQTask)
	}
	internal := r.Group("/internal")
	{
		internal.PUT("/quantum_tasks/updating", handler.SetFraxResult)
	}
}

// RegisterStatic регистрирует статические файлы и шаблоны
// @Description Настраивает обслуживание статических файлов и HTML шаблонов
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler - внутренний вспомогательный метод для обработки ошибок
// Не экспортируется в Swagger документацию
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":  "fail",
		"message": err.Error(),
	})
}

// okJSON - внутренний вспомогательный метод для успешных ответов
// Не экспортируется в Swagger документацию
func (h *Handler) okJSON(ctx *gin.Context, statusCode int, payload interface{}) {
	ctx.JSON(statusCode, gin.H{
		"status": "ok",
		"data":   payload,
	})
}
