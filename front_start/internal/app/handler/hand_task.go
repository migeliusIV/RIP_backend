package handler

import (
    "errors"
    "front_start/internal/app/ds"
    "net/http"
    "strconv"
    "time"

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
			ID_user:      hardcodedUserID,
			TaskStatus:   ds.StatusDraft,
			CreationDate: time.Now(),
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

	c.HTML(http.StatusOK, "quantum_task.html", task)
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

type taskUpdateRequest struct {
    TaskDescription string  `json:"task_description"`
}

func (h *Handler) ApiListTasks(ctx *gin.Context) {
    status := ctx.Query("status")
    from := ctx.Query("from")
    to := ctx.Query("to")
    tasks, err := h.Repository.ListTasks(status, from, to)
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, gin.H{"items": tasks})
}

func (h *Handler) ApiGetTask(ctx *gin.Context) {
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

func (h *Handler) ApiUpdateTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil || id <= 0 {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    var req taskUpdateRequest
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

func (h *Handler) ApiFormTask(ctx *gin.Context) {
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

type resolveRequest struct {
    Action string `json:"action"` // "complete" | "reject"
}

func (h *Handler) ApiResolveTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil || id <= 0 {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    var req resolveRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    resolved, err := h.Repository.ResolveTask(uint(id), hardcodedUserID, req.Action, time.Now())
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, resolved)
}

func (h *Handler) ApiDeleteTask(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil || id <= 0 {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    if err := h.Repository.DeleteTask(uint(id)); err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, gin.H{"id": id})
}

// ---- JSON API for m-m ----
type updateDegreesRequest struct {
    Degrees *float32 `json:"degrees"`
}

func (h *Handler) ApiRemoveServiceFromTask(ctx *gin.Context) {
    taskID, err1 := strconv.Atoi(ctx.Param("task_id"))
    gateID, err2 := strconv.Atoi(ctx.Param("service_id"))
    if err1 != nil || err2 != nil || taskID <= 0 || gateID <= 0 {
        h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid ids"))
        return
    }
    if err := h.Repository.RemoveServiceFromTask(uint(taskID), uint(gateID)); err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, gin.H{"task_id": taskID, "service_id": gateID})
}

func (h *Handler) ApiUpdateDegrees(ctx *gin.Context) {
    taskID, err1 := strconv.Atoi(ctx.Param("task_id"))
    gateID, err2 := strconv.Atoi(ctx.Param("service_id"))
    if err1 != nil || err2 != nil || taskID <= 0 || gateID <= 0 {
        h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid ids"))
        return
    }
    var req updateDegreesRequest
    if err := ctx.ShouldBindJSON(&req); err != nil || req.Degrees == nil {
        h.errorHandler(ctx, http.StatusBadRequest, errors.New("degrees required"))
        return
    }
    if err := h.Repository.UpdateDegrees(uint(taskID), uint(gateID), *req.Degrees); err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, gin.H{"task_id": taskID, "service_id": gateID, "degrees": req.Degrees})
}
