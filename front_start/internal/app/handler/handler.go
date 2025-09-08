package handler

import (
	"front_start/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetGates(ctx *gin.Context) {
	var gates []repository.Gate
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		gates, err = h.Repository.GetGates()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		gates, err = h.Repository.GetGatesByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"gates": gates,
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetGate(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	gate, err := h.Repository.GetGate(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "properties.html", gin.H{
		"gate": gate,
	})
}
