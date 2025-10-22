package handler

import (
    "front_start/internal/app/config"
    appredis "front_start/internal/app/redis"
    "front_start/internal/app/repository"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
    Redis      *appredis.Client
	JWTConfig  *config.JWTConfig
}

func NewHandler(r *repository.Repository, redis *appredis.Client, jwtConfig *config.JWTConfig) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redis,
		JWTConfig:  jwtConfig,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (handler *Handler) RegisterHandler(r *gin.Engine) {
	r.POST("/login", handler.Login)
	r.POST("/users", handler.Register)
	r.GET("/IBM", handler.GetGates)

	// Эндпоинты, доступные только модераторам
    moderator := r.Group("/")
    moderator.Use(handler.AuthMiddleware, handler.ModeratorMiddleware)
	{
		// Управление факторами (создание, изменение, удаление)
		moderator.POST("/api/gates", handler.ApiAddGate)
		moderator.PUT("/api/gates/:id", handler.ApiUpdateGate)
		moderator.DELETE("/api/gates/:id", handler.ApiDeleteGate)
		moderator.POST("/api/gates/:id/image", handler.ApiUploadGatesImage)

		// Управление заявками (завершение/отклонение)
		moderator.PUT("/api/quantum_tasks/:id/resolve", handler.ApiResolveQTask)
	}
	// Эндпоинты, доступные всем авторизованным пользователям
    auth := r.Group("/")
    auth.Use(handler.AuthMiddleware)
	{
		// Пользователи
		auth.POST("/api/auth/logout", handler.Logout)
		auth.GET("/api/users/me", handler.ApiMe)
		auth.PUT("/api/users/me", handler.ApiUpdateMe)
		// m-m (2)
		auth.DELETE("/api/tasks/:task_id/services/:service_id", handler.ApiRemoveGateFromTask)
		auth.PUT("/api/tasks/:task_id/services/:service_id", handler.ApiUpdateDegrees)
		// каша
		//auth.GET("/IBM", handler.GetGates)
		auth.GET("/gate_property/:id", handler.GetGateByID)
		auth.GET("/quantum_task/:id", handler.GetTask)
		auth.POST("/quantum_task/add/gate/:id_gate", handler.AddGateToTask) // orm
		auth.POST("/quantum_task/:task_id/delete", handler.DeleteTask)      // удаление заявки через SQL
		// API
		// Gates (7)
		auth.GET("/api/gates", handler.ApiGatesList)
		auth.GET("/api/gates/:id", handler.ApiGetGateByID)
		auth.POST("/api/draft/gates/:id", handler.ApiAddGateToDraft)

		// Quantum tasks (7)
		auth.GET("/api/quantum_task/current", handler.ApiGetCurrQTask)
		auth.GET("/api/quantum_tasks", handler.ApiListQTasks)
		auth.GET("/api/quantum_tasks/:id", handler.ApiGetQTaskByID)
		auth.PUT("/api/quantum_tasks/:id", handler.ApiUpdateQTask)
		auth.PUT("/api/quantum_tasks/:id/form", handler.ApiFormQTask)
		auth.DELETE("/api/quantum_tasks/:id", handler.ApiDeleteQTask)
	}
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":  "fail",
		"message": err.Error(),
	})
}

// okJSON отправляет успешный JSON ответ с произвольным payload
func (h *Handler) okJSON(ctx *gin.Context, statusCode int, payload interface{}) {
	ctx.JSON(statusCode, gin.H{
		"status": "ok",
		"data":   payload,
	})
}
