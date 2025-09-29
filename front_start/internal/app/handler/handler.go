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
	//r.POST("/quantum_task/add/gate/:id_gate", handler.AddGateToTask) // orm
	//r.POST("/quantum_task/:task_id/delete", handler.DeleteTask)      // удаление заявки через SQL

	// Домен услуг (гейтов)
	r.GET("/IBM", handler.GetGates)
	r.GET("/gate_property/:id", handler.GetGateByID)
	// r.POST("/factors", h.CreateFactor)
	// r.PUT("/factors/:id", h.UpdateFactor)
	// r.DELETE("/factors/:id", h.DeleteFactor)
	// r.POST("/frax/draft/factors/:factor_id", h.AddFactorToDraft)
	// r.POST("/factors/:id/image", h.UploadFactorImage)

	// Домен заявок (задач)
	r.GET("/quantum_task/:id", handler.GetTask)
	// r.GET("/frax/cart", h.GetCartBadge)
	// r.GET("/frax", h.ListFrax)
	// r.PUT("/frax/:id", h.UpdateFrax)
	// r.PUT("/frax/:id/form", h.FormFrax)
	// r.PUT("/frax/:id/resolve", h.ResolveFrax)
	// r.DELETE("/frax/:id", h.DeleteFrax)

	// Домен м-м
	// r.DELETE("/frax/:id/factors/:factor_id", h.RemoveFactorFromFrax)
	// r.PUT("/frax/:id/factors/:factor_id", h.UpdateMM)

	// Домен пользователей
	// r.POST("/users", h.Register)
	// r.GET("/users/:id", h.GetUserData)
	// r.PUT("/users/:id", h.UpdateUserData)
	// r.POST("/auth/login", h.Login)
	// r.POST("/auth/logout", h.Logout)
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
		"status":      "error",
		"description": err.Error(),
	})
}
