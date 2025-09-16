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
	r.GET("/task/:id", handler.GetTask)
	//r.POST("/task/add/gate/:gate_id", handler.) - добавление в заявку через ORM
	//r.POST("/task/:task_id/delete", handler.) - удаление заявки через SQL
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
