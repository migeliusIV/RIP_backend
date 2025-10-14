package handler

import (
	"errors"
	"fmt"
	"front_start/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const hardcodedUserID = 1

func (h *Handler) AddGateToTask(c *gin.Context) {
	gateID, err := strconv.Atoi(c.Param("id_gate"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}

	task, err := h.Repository.GetDraftTask(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newTask := ds.QuantumTask{
			ID_user:      userID,
			TaskStatus:   ds.StatusDraft,
			CreationDate: time.Now(),
			Res_koeff_0:  1.0,  // Начальное значение для |0⟩
			Res_koeff_1:  0.0,  // Начальное значение для |1⟩
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
		// Возвращаем кастомную 404-страницу
		c.HTML(http.StatusNotFound, "invalid_taskpage.html", nil)
		return
	}

	if len(task.GatesDegrees) == 0 {
		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty frax page, add factors first"))
		return
	}
	
	var taskRepresent ds.QuantumTask
	taskRepresent = *task  // разыменовываем указатель
	taskRepresent.Res_koeff_0 = task.Res_koeff_0 * task.Res_koeff_0  // квадрат первого коэффициента
	taskRepresent.Res_koeff_1 = task.Res_koeff_1 * task.Res_koeff_1  // квадрат второго коэффициента
	c.HTML(http.StatusOK, "quantum_task.html", taskRepresent)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("task_id"))

	if err := h.Repository.LogicallyDeleteTask(uint(taskID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/IBM")
}

// ---- JSON API for tasks ----

// moved to serialization.go

func (h *Handler) ApiListQTasks(ctx *gin.Context) {
	status := ctx.Query("status")
	from := ctx.Query("from")
	to := ctx.Query("to")
	tasks, err := h.Repository.ListTasks(status, from, to)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
    h.okJSON(ctx, http.StatusOK, DTO_TasksListResponse{Items: tasks})
}

func (h *Handler) ApiGetQTaskByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	task, err := h.Repository.GetTaskWithGates(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	h.okJSON(ctx, http.StatusOK, task)
}

func (h *Handler) ApiUpdateQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
    var req DTO_taskUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	updated, err := h.Repository.UpdateTask(uint(id), req.TaskDescription)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.okJSON(ctx, http.StatusOK, updated)
}

func (h *Handler) ApiFormQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	formed, err := h.Repository.FormTask(uint(id), time.Now())
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	h.okJSON(ctx, http.StatusOK, formed)
}

// moved to serialization.go

func (h *Handler) ApiResolveQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		logrus.Errorf("Invalid task ID in resolve request: %v", err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	
	logrus.Infof("Processing resolve request for task ID: %d", id)
	
    var req DTO_resolveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Более детальная обработка ошибки JSON
		logrus.Errorf("JSON binding error for task %d: %v", id, err)
		if err.Error() == "EOF" {
			h.errorHandler(ctx, http.StatusBadRequest, errors.New("request body is empty, expected JSON with 'action' field"))
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid JSON format: %v", err))
		}
		return
	}
	
	// Проверяем, что action указан
	if req.Action == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("action field is required"))
		return
	}
	
	// Проверяем, что action имеет допустимое значение
	if req.Action != "complete" && req.Action != "reject" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("action must be 'complete' or 'reject'"))
		return
	}
	
	// Если задача завершается, вычисляем результат
	if req.Action == "complete" {
		// Проверяем, что у задачи есть все необходимые параметры для вычисления результата
		task, err := h.Repository.GetTaskWithGates(uint(id))
		if err != nil {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}
		
		// Проверяем наличие описания задачи
		if task.TaskDescription == "" {
			h.errorHandler(ctx, http.StatusBadRequest, errors.New("task description is required for completion"))
			return
		}
		
		// Проверяем, что у всех гейтов указаны градусы
		for _, gateDegree := range task.GatesDegrees {
			var currGate = gateDegree.Gate

			if gateDegree.Degrees == nil && currGate.TheAxis != "non"{
				h.errorHandler(ctx, http.StatusBadRequest, 
					fmt.Errorf("degrees not specified for gate %s (ID: %d)", 
						gateDegree.Gate.Title, gateDegree.Gate.ID_gate))
				return
			}
		}
		
		// Вычисляем результат квантовой задачи
		logrus.Infof("Calculating result for task %d", id)
		err = h.Repository.GetQTaskRes(uint(id))
		if err != nil {
			logrus.Errorf("Failed to calculate result for task %d: %v", id, err)
			h.errorHandler(ctx, http.StatusInternalServerError, 
				fmt.Errorf("failed to calculate task result: %v", err))
			return
		}
		logrus.Infof("Successfully calculated result for task %d", id)
	}
	
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}
	adminID := uint(userID)
	resolved, err := h.Repository.ResolveTask(uint(id), adminID, req.Action, time.Now())
	if err != nil {
		logrus.Errorf("Failed to resolve task %d: %v", id, err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	
	logrus.Infof("Successfully resolved task %d with action: %s", id, req.Action)
	h.okJSON(ctx, http.StatusOK, resolved)
}

func (h *Handler) ApiDeleteQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.DeleteTask(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
    h.okJSON(ctx, http.StatusOK, DTO_SimpleIDResponse{ID: id})
}

// ---- JSON API for m-m ----
// moved to serialization.go

func (h *Handler) ApiRemoveGateFromTask(ctx *gin.Context) {
	taskID, err1 := strconv.Atoi(ctx.Param("task_id"))
	gateID, err2 := strconv.Atoi(ctx.Param("service_id"))
	if err1 != nil || err2 != nil || taskID <= 0 || gateID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid ids"))
		return
	}
	if err := h.Repository.RemoveGateFromTask(uint(taskID), uint(gateID)); err != nil { //RemoveServiceFromTask -> RemoveGateFromTask
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
    h.okJSON(ctx, http.StatusOK, DTO_TaskServiceLinkResponse{TaskID: uint(taskID), ServiceID: gateID})
}

func (h *Handler) ApiUpdateDegrees(ctx *gin.Context) {
	taskID, err1 := strconv.Atoi(ctx.Param("task_id"))
	gateID, err2 := strconv.Atoi(ctx.Param("service_id"))
	if err1 != nil || err2 != nil || taskID <= 0 || gateID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid ids"))
		return
	}
    var req DTO_updateDegreesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Degrees == nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("degrees required"))
		return
	}
	if err := h.Repository.UpdateDegrees(uint(taskID), uint(gateID), *req.Degrees); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
    h.okJSON(ctx, http.StatusOK, DTO_UpdateDegreesResponse{TaskID: taskID, ServiceID: gateID, Degrees: req.Degrees})
}
