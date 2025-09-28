package handler

import (
	"errors"
	"front_start/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const hardcodedUserID = 1

func (h *Handler) AddGateToTask(c *gin.Context) {
	gateID, err := strconv.Atoi(c.Param("id_gate"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	task, err := h.Repository.GetDraftTask(hardcodedUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newTask := ds.QuantumTask{
			ID_user:    hardcodedUserID,
			TaskStatus: ds.StatusDraft,
		}
		if createErr := h.Repository.CreateTask(&newTask); createErr != nil {
			h.errorHandler(c, http.StatusInternalServerError, createErr)
			return
		}
		task = &newTask
	} else if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if err = h.Repository.AddGateToTask(task.ID_task, uint(gateID)); err != nil {
	}

	c.Redirect(http.StatusFound, "/IBM")
}

func (h *Handler) GetTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	task, err := h.Repository.GetTaskWithGates(uint(taskID))
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	if len(task.GatesDegrees) == 0 {
		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty frax page, add factors first"))
		return
	}

	c.HTML(http.StatusOK, "quantum_task.html", task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("id_task"))

	if err := h.Repository.LogicallyDeleteFrax(uint(taskID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/IBM")
}
