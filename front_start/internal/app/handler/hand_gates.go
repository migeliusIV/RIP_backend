package handler

import (
	"front_start/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetGates(ctx *gin.Context) {
	var gates []ds.Gate
	var err error

	search := ctx.Query("gateSearching")
	if search == "" {
		gates, err = h.Repository.GetGates()
	} else {
		gates, err = h.Repository.GetGatesByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	draftTask, _ := h.Repository.GetDraftTask(hardcodedUserID)
	var taskID uint = 0
	var gatesCount int = 0

	if draftTask != nil {
		// Загружаем связанные факторы, чтобы посчитать их количество
		fullTask, err := h.Repository.GetTaskWithGates(draftTask.ID_task)
		if err == nil {
			taskID = fullTask.ID_task
			gatesCount = len(fullTask.Task)
		}
	}

	ctx.HTML(http.StatusOK, "gates_list.html", gin.H{
		"gates":         gates,
		"gateSearching": search,
		"taskID":        taskID,
		"gatesCount":    gatesCount,
	})
}

func (h *Handler) GetGateByID(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	gate, err := h.Repository.GetGateByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "gate_properties.html", gate)
}
