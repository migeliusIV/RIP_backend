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
			Res_koeff_0:  1.0, // Начальное значение для |0⟩
			Res_koeff_1:  0.0, // Начальное значение для |1⟩
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
	taskRepresent = *task                                           // разыменовываем указатель
	taskRepresent.Res_koeff_0 = task.Res_koeff_0 * task.Res_koeff_0 // квадрат первого коэффициента
	taskRepresent.Res_koeff_1 = task.Res_koeff_1 * task.Res_koeff_1 // квадрат второго коэффициента
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

// ApiListQTasks возвращает список задач
// @Summary Получить список задач
// @Description Возвращает список квантовых задач. Пользователь видит только свои задачи, модератор - все задачи, неавторизованный пользователь - ошибку доступа.
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу"
// @Param from query string false "Начальная дата (фильтр от)"
// @Param to query string false "Конечная дата (фильтр до)"
// @Success 200 {array} DTO_Resp_Tasks "Список задач"
// @Failure 401 {object} string "Unauthorized"
// @Failure 403 {object} string "Forbidden"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /api/quantum_tasks [get]
func (h *Handler) ApiListQTasks(ctx *gin.Context) {
	// Получаем ID пользователя из контекста (устанавливается в middleware аутентификации)
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("требуется авторизация"))
		return
	}

	// Получаем информацию о пользователе
	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	status := ctx.Query("status")
	from := ctx.Query("from")
	to := ctx.Query("to")

	var tasks []*ds.QuantumTask

	// Разделяем логику в зависимости от роли пользователя
	if user.IsAdmin {
		// Модератор видит все задачи
		tasks, err = h.Repository.ListTasks(status, from, to)
	} else {
		// Обычный пользователь видит только свои задачи
		tasks, err = h.Repository.ListTasksByUser(userID, status, from, to)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Преобразуем задачи в DTO
	var represent_tasks []DTO_Resp_Tasks
	for _, task := range tasks {
		// 1. Преобразуем GatesDegrees
		var dtoGatesDegrees []DTO_Resp_GatesDegrees
		for _, gateDegree := range task.GatesDegrees {
			dtoGatesDegrees = append(dtoGatesDegrees, DTO_Resp_GatesDegrees{
				ID_gate: gateDegree.ID_gate,
				ID_task: gateDegree.ID_task,
				Degrees: gateDegree.Degrees,
			})
		}

		// 2. Создаём DTO задачи
		represent_tasks = append(represent_tasks, DTO_Resp_Tasks{
			ID_task:         task.ID_task,
			TaskStatus:      task.TaskStatus,
			CreationDate:    task.CreationDate,
			ID_user:         task.ID_user,
			ConclusionDate:  task.ConclusionDate,
			FormedDate:      task.FormedDate,
			TaskDescription: task.TaskDescription,
			Res_koeff_0:     task.Res_koeff_0,
			Res_koeff_1:     task.Res_koeff_1,
			GatesDegrees:    dtoGatesDegrees,
		})
	}

	ctx.JSON(http.StatusOK, represent_tasks)
}

// ApiGetQTaskByID возвращает задачу по ID
// @Summary Получить задачу по ID
// @Description Возвращает детальную информацию о квантовой задаче по её идентификатору
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} DTO_Resp_Tasks "Детали задачи"
// @Failure 400 {object} string "Invalid task ID"
// @Failure 404 {object} string "Task not found"
// @Router /api/quantum_tasks/{id} [get]
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

	var represent_task DTO_Resp_Tasks
	var dtoGatesDegrees []DTO_Resp_GatesDegrees
	for _, gateDegree := range task.GatesDegrees {
		gateInfo, err := h.Repository.GetGateByID(int(gateDegree.ID_gate))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}

		dtoGatesDegrees = append(dtoGatesDegrees, DTO_Resp_GatesDegrees{
			Title:   gateInfo.Title, // или gateDegree.Gate.ID_gate, если нужно
			TheAxis: gateInfo.TheAxis,
			Image:   gateInfo.Image,
			ID_gate: gateInfo.ID_gate,
			ID_task: gateDegree.ID_task, // или gateDegree.Task.ID_task
			Degrees: gateDegree.Degrees,
		})
	}

	represent_task = DTO_Resp_Tasks{
		ID_task:         task.ID_task,
		TaskStatus:      task.TaskStatus,
		CreationDate:    task.CreationDate,
		ID_user:         task.ID_user,
		ConclusionDate:  task.ConclusionDate,
		TaskDescription: task.TaskDescription,
		Res_koeff_0:     task.Res_koeff_0,
		Res_koeff_1:     task.Res_koeff_1,
		GatesDegrees:    dtoGatesDegrees,
	}
	ctx.JSON(http.StatusOK, represent_task)
}

// ApiUpdateQTask обновляет описание задачи
// @Summary Обновить задачу
// @Description Обновляет описание квантовой задачи
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body DTO_Req_TaskUpd true "Данные для обновления"
// @Success 200 {object} DTO_Resp_Tasks "Обновленная задача"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Router /api/quantum_tasks/{id} [put]
func (h *Handler) ApiUpdateQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var req DTO_Req_TaskUpd
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	task, err := h.Repository.UpdateTask(uint(id), req.TaskDescription)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	var represent_task DTO_Resp_Tasks
	var dtoGatesDegrees []DTO_Resp_GatesDegrees
	for _, gateDegree := range task.GatesDegrees {
		dtoGatesDegrees = append(dtoGatesDegrees, DTO_Resp_GatesDegrees{
			ID_gate: gateDegree.ID_gate, // или gateDegree.Gate.ID_gate, если нужно
			ID_task: gateDegree.ID_task, // или gateDegree.Task.ID_task
			Degrees: gateDegree.Degrees,
		})
	}

	represent_task = DTO_Resp_Tasks{
		ID_task:         task.ID_task,
		TaskStatus:      task.TaskStatus,
		CreationDate:    task.CreationDate,
		ID_user:         task.ID_user,
		ConclusionDate:  task.ConclusionDate,
		TaskDescription: task.TaskDescription,
		Res_koeff_0:     task.Res_koeff_0,
		Res_koeff_1:     task.Res_koeff_1,
		GatesDegrees:    dtoGatesDegrees,
	}
	ctx.JSON(http.StatusOK, represent_task)
}

// ApiFormQTask формирует задачу
// @Summary Сформировать задачу
// @Description Переводит задачу из статуса черновика в статус сформированной
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} DTO_Resp_Tasks "Сформированная задача"
// @Failure 400 {object} string "Invalid task ID or cannot form task"
// @Router /api/quantum_tasks/{id}/form [put]
func (h *Handler) ApiFormQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	task, err := h.Repository.FormTask(uint(id), time.Now())
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var represent_task DTO_Resp_Tasks
	var dtoGatesDegrees []DTO_Resp_GatesDegrees
	for _, gateDegree := range task.GatesDegrees {
		dtoGatesDegrees = append(dtoGatesDegrees, DTO_Resp_GatesDegrees{
			ID_gate: gateDegree.ID_gate, // или gateDegree.Gate.ID_gate, если нужно
			ID_task: gateDegree.ID_task, // или gateDegree.Task.ID_task
			Degrees: gateDegree.Degrees,
		})
	}

	represent_task = DTO_Resp_Tasks{
		ID_task:         task.ID_task,
		TaskStatus:      task.TaskStatus,
		CreationDate:    task.CreationDate,
		ID_user:         task.ID_user,
		ConclusionDate:  task.ConclusionDate,
		TaskDescription: task.TaskDescription,
		Res_koeff_0:     task.Res_koeff_0,
		Res_koeff_1:     task.Res_koeff_1,
		GatesDegrees:    dtoGatesDegrees,
	}
	ctx.JSON(http.StatusOK, represent_task)
}

// ApiResolveQTask завершает или отклоняет задачу
// @Summary Завершить/отклонить задачу
// @Description Выполняет завершение или отклонение квантовой задачи с вычислением результата
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param request body DTO_Req_TaskResolve true "Действие с задачей"
// @Success 200 {object} DTO_Resp_Tasks "Решенная задача"
// @Failure 400 {object} string "Invalid input or missing required fields"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/quantum_tasks/{id}/resolve [put]
func (h *Handler) ApiResolveQTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		logrus.Errorf("Invalid task ID in resolve request: %v", err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Processing resolve request for task ID: %d", id)

	var req DTO_Req_TaskResolve
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

			if gateDegree.Degrees == nil && currGate.TheAxis != "non" {
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

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	adminID := uint(userID)
	task, err := h.Repository.ResolveTask(uint(id), adminID, req.Action, time.Now())
	if err != nil {
		logrus.Errorf("Failed to resolve task %d: %v", id, err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var represent_task DTO_Resp_Tasks
	var dtoGatesDegrees []DTO_Resp_GatesDegrees
	for _, gateDegree := range task.GatesDegrees {
		dtoGatesDegrees = append(dtoGatesDegrees, DTO_Resp_GatesDegrees{
			ID_gate: gateDegree.ID_gate, // или gateDegree.Gate.ID_gate, если нужно
			ID_task: gateDegree.ID_task, // или gateDegree.Task.ID_task
			Degrees: gateDegree.Degrees,
		})
	}

	represent_task = DTO_Resp_Tasks{
		ID_task:         task.ID_task,
		TaskStatus:      task.TaskStatus,
		CreationDate:    task.CreationDate,
		ID_user:         task.ID_user,
		ConclusionDate:  task.ConclusionDate,
		TaskDescription: task.TaskDescription,
		Res_koeff_0:     task.Res_koeff_0,
		Res_koeff_1:     task.Res_koeff_1,
		GatesDegrees:    dtoGatesDegrees,
	}

	logrus.Infof("Successfully resolved task %d with action: %s", id, req.Action)
	ctx.JSON(http.StatusOK, represent_task)
}

// ApiDeleteQTask удаляет задачу
// @Summary Удалить задачу
// @Description Полностью удаляет квантовую задачу из системы
// @Tags QuantumTasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} DTO_Resp_SimpleID "ID удаленной задачи"
// @Failure 400 {object} string "Invalid task ID"
// @Failure 500 {object} string "Internal server error"
// @Router /api/quantum_tasks/{id} [delete]
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
	ctx.JSON(http.StatusOK, DTO_Resp_SimpleID{ID: id})
}

// ---- JSON API for m-m ----
// moved to serialization.go

// ApiRemoveGateFromTask удаляет гейт из задачи
// @Summary Удалить гейт из задачи
// @Description Удаляет связь между гейтом и задачей
// @Tags M-M
// @Accept json
// @Produce json
// @Param task_id path int true "ID задачи"
// @Param service_id path int true "ID гейта"
// @Success 200 {object} DTO_Resp_TaskServiceLink "Информация об удаленной связи"
// @Failure 400 {object} string "Invalid IDs"
// @Failure 500 {object} string "Internal server error"
// @Router /api/tasks/{task_id}/services/{service_id} [delete]
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
	ctx.JSON(http.StatusOK, DTO_Resp_TaskServiceLink{TaskID: uint(taskID), ServiceID: gateID})
}

// ApiUpdateDegrees обновляет градусы для гейта в задаче
// @Summary Обновить градусы гейта
// @Description Обновляет значение градусов для конкретного гейта в задаче
// @Tags M-M
// @Accept json
// @Produce json
// @Param task_id path int true "ID задачи"
// @Param service_id path int true "ID гейта"
// @Param request body DTO_Req_DegreesUpd true "Новые значения градусов"
// @Success 200 {object} DTO_Resp_UpdateDegrees "Обновленные данные"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Router /api/tasks/{task_id}/services/{service_id} [put]
func (h *Handler) ApiUpdateDegrees(ctx *gin.Context) {
	taskID, err1 := strconv.Atoi(ctx.Param("task_id"))
	gateID, err2 := strconv.Atoi(ctx.Param("service_id"))
	if err1 != nil || err2 != nil || taskID <= 0 || gateID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid ids"))
		return
	}
	var req DTO_Req_DegreesUpd
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Degrees == nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("degrees required"))
		return
	}
	if err := h.Repository.UpdateDegrees(uint(taskID), uint(gateID), *req.Degrees); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, DTO_Resp_UpdateDegrees{TaskID: taskID, ServiceID: gateID, Degrees: req.Degrees})
}

// PUT /api/internal/frax/result
func (h *Handler) SetFraxResult(c *gin.Context) {
	token := c.GetHeader("Authorization")
	expectedToken := "secret12"

	if token != expectedToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}

	var res DTO_Res_TaskUpd
	if err := c.ShouldBindJSON(&res); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	task, err := h.Repository.UpdateTask(uint(res.ID_task), res.TaskDescription)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, task)
}
