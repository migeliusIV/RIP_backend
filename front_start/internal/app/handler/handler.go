package handler

import (
	"front_start/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (handler *Handler) RegisterHandler(r *gin.Engine) {
	r.GET("/IBM", handler.GetGates)
	r.GET("/gate_property/:id", handler.GetGateByID)
	r.GET("/quantum_task/:id", handler.GetTask)
	r.POST("/quantum_task/add/gate/:id_gate", handler.AddGateToTask) // orm
	r.POST("/quantum_task/:task_id/delete", handler.DeleteTask)      // удаление заявки через SQL
    // JSON API routes (exactly 21)
    // Services (7)
    r.GET("/api/gates", handler.ApiListServices)
    r.GET("/api/gates/:id", handler.ApiGetService)
    r.POST("/api/gates", handler.ApiCreateService)
    r.PUT("/api/gates/:id", handler.ApiUpdateService)
    r.DELETE("/api/gates/:id", handler.ApiDeleteService)
    r.POST("/api/draft/gates/:id", handler.ApiAddServiceToDraft)
    r.POST("/api/gates/:id/image", handler.ApiUploadServiceImage)

    // Tasks (7)
    r.GET("/api/cart", handler.ApiGetCartBadge)
    r.GET("/api/quantum_tasks", handler.ApiListTasks)
    r.GET("/api/quantum_tasks/:id", handler.ApiGetTask)
    r.PUT("/api/quantum_tasks/:id", handler.ApiUpdateTask)
    r.PUT("/api/quantum_tasks/:id/form", handler.ApiFormTask)
    r.PUT("/api/quantum_tasks/:id/resolve", handler.ApiResolveTask)
    r.DELETE("/api/quantum_tasks/:id", handler.ApiDeleteTask)

    // m-m (2)
    r.DELETE("/api/tasks/:task_id/services/:service_id", handler.ApiRemoveServiceFromTask)
    r.PUT("/api/tasks/:task_id/services/:service_id", handler.ApiUpdateDegrees)

    // Users (5)
    r.POST("/api/users/register", handler.ApiRegister)
    r.GET("/api/users/me", handler.ApiMe)
    r.PUT("/api/users/me", handler.ApiUpdateMe)
    r.POST("/api/auth/login", handler.ApiLogin)
    r.POST("/api/auth/logout", handler.ApiLogout)
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
