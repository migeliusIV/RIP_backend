package handler

import (
	"front_start/internal/app/ds"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"errors"
	
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
	
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("требуется авторизация"))
		return
	}

	draftTask, _ := h.Repository.GetDraftTask(userID)
	var taskID uint = 0
	var gatesCount int = 0

	if draftTask != nil {
		// Загружаем связанные факторы, чтобы посчитать их количество
		fullTask, err := h.Repository.GetTaskWithGates(draftTask.ID_task)
		if err == nil {
			taskID = fullTask.ID_task
			gatesCount = len(fullTask.GatesDegrees)
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

// ---- JSON API (services/gates) ----

// ApiGatesList godoc
// @Summary      Получить список гейтов
// @Description  Возвращает список всех гейтов. Поддерживает фильтрацию по названию.
// @Tags         Gates
// @Produce      json
// @Param        title query string false "Фильтр по названию гейта (поиск по подстроке)"
// @Success      200 {array} DTO_Resp_Gate
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /api/gates [get]
func (h *Handler) ApiGatesList(ctx *gin.Context) {
	title := ctx.Query("title")
	gates, err := h.Repository.ListGates(title)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	var represent_gates []DTO_Resp_Gate
	for _, gate := range gates {
		represent_gates = append(represent_gates, DTO_Resp_Gate{
			ID_gate:      gate.ID_gate,
			Title:        gate.Title,
			Description:  gate.Description,
			Status:       gate.Status,
			Image:        gate.Image,
			I0j0:         gate.I0j0,
			I0j1:         gate.I0j1,
			I1j0:         gate.I1j0,
			I1j1:         gate.I1j1,
			Matrix_koeff: gate.Matrix_koeff,
			// subject area
			FullInfo: gate.FullInfo,
			TheAxis:  gate.TheAxis,
		})
	}
	ctx.JSON(http.StatusOK, represent_gates)
}

// ApiGetGateByID godoc
// @Summary      Получить гейт по ID
// @Description  Возвращает полную информацию о гейте по его идентификатору.
// @Tags         Gates
// @Produce      json
// @Param        id path int true "ID гейта"
// @Success      200 {object} ds.Gate
// @Failure      400 {object} map[string]string "Некорректный ID"
// @Failure      404 {object} map[string]string "Гейт не найден"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /api/gates/{id} [get]
func (h *Handler) ApiGetGateByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	gate, err := h.Repository.GetGateByID(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, gate)
}

// ApiAddGate godoc
// @Summary      Создать новый гейт
// @Description  Создаёт новый квантовый гейт с указанными параметрами.
// @Tags         Gates
// @Accept       json
// @Produce      json
// @Param        gate body DTO_Req_GateCreate true "Данные нового гейта"
// @Success      201 {object} ds.Gate
// @Failure      400 {object} map[string]string "Некорректные данные запроса"
// @Failure      500 {object} map[string]string "Ошибка при создании гейта"
// @Router       /api/gates [post]
func (h *Handler) ApiAddGate(ctx *gin.Context) {
	var req DTO_Req_GateCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if req.Title == "" || req.Description == "" || req.FullInfo == "" {
		h.errorHandler(ctx, http.StatusBadRequest, gin.Error{Err: gin.Error{}})
		return
	}
	gate := ds.Gate{
		Title:        req.Title,
		Description:  req.Description,
		FullInfo:     req.FullInfo,
		TheAxis:      req.TheAxis,
		I0j0:         req.I0j0,
		I0j1:         req.I0j1,
		I1j0:         req.I1j0,
		I1j1:         req.I1j1,
		Matrix_koeff: req.Matrix_koeff,
	}
	if req.Status != nil {
		gate.Status = *req.Status
	}
	if err := h.Repository.AddGate(&gate); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, gate)
}

// ApiUpdateGate godoc
// @Summary      Обновить гейт
// @Description  Обновляет существующий гейт по ID.
// @Tags         Gates
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID гейта"
// @Param        gate body DTO_Req_GateCreate true "Обновлённые данные гейта"
// @Success      200 {object} ds.Gate
// @Failure      400 {object} map[string]string "Некорректные данные запроса или ID"
// @Failure      500 {object} map[string]string "Ошибка при обновлении гейта"
// @Router       /api/gates/{id} [put]
func (h *Handler) ApiUpdateGate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var req DTO_Req_GateCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	//UpdateService -> UpdateGate
	updated, err := h.Repository.UpdateGate(uint(id), req.Title, req.Description, req.FullInfo, req.TheAxis, req.Status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// ApiDeleteGate godoc
// @Summary      Удалить гейт
// @Description  Удаляет гейт по ID.
// @Tags         Gates
// @Produce      json
// @Param        id path int true "ID гейта"
// @Success      200 {object} DTO_Resp_SimpleID
// @Failure      400 {object} map[string]string "Некорректный ID"
// @Failure      500 {object} map[string]string "Ошибка при удалении гейта"
// @Router       /api/gates/{id} [delete]
func (h *Handler) ApiDeleteGate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	//DeleteService -> DeleteGate
	if err := h.Repository.DeleteGate(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, DTO_Resp_SimpleID{ID: id})
}

// ApiAddGateToDraft godoc
// @Summary      Добавить гейт в черновик задачи
// @Description  Добавляет указанный гейт в текущую черновую задачу пользователя.
// @Tags         Gates
// @Produce      json
// @Param        id path int true "ID гейта"
// @Success      201 {object} DTO_Resp_TaskServiceLink
// @Failure      400 {object} map[string]string "Некорректный ID гейта"
// @Failure      401 {object} map[string]string "Требуется авторизация"
// @Failure      500 {object} map[string]string "Ошибка при добавлении гейта в задачу"
// @Router       /api/draft/gates/{id} [post]
func (h *Handler) ApiAddGateToDraft(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	// Reuse HTML flow: get or create draft, then add gate to task
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	task, err := h.Repository.GetDraftTask(userID)
	if err != nil {
		// create
		newTask := ds.QuantumTask{
			ID_user:      userID,
			TaskStatus:   ds.StatusDraft,
			CreationDate: time.Now(),
		}
		if createErr := h.Repository.CreateTask(&newTask); createErr != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, createErr)
			return
		}
		task = &newTask
	}
	if err := h.Repository.AddGateToTask(task.ID_task, uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, DTO_Resp_TaskServiceLink{TaskID: task.ID_task, ServiceID: id})
}

// ApiGetCurrQTask godoc
// @Summary      Получить информацию о текущей черновой задаче
// @Description  Возвращает ID текущей черновой задачи и количество добавленных в неё гейтов.
// @Tags         QuantumTasks
// @Produce      json
// @Success      200 {object} DTO_Resp_CurrTaskInfo
// @Failure      500 {object} map[string]string "Ошибка при получении данных задачи"
// @Router       /api/quantum_task/current [get]
func (h *Handler) ApiGetCurrQTask(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("требуется авторизация"))
		return
	}
	draftTask, _ := h.Repository.GetDraftTask(userID)
	var taskID uint = 0
	var gatesCount int = 0
	if draftTask != nil {
		fullTask, err := h.Repository.GetTaskWithGates(draftTask.ID_task)
		if err == nil {
			taskID = fullTask.ID_task
			gatesCount = len(fullTask.GatesDegrees)
		}
	}
	ctx.JSON(http.StatusOK, DTO_Resp_CurrTaskInfo{TaskID: taskID, ServicesCount: gatesCount})
}

// ApiUploadGatesImage godoc
// @Summary      Загрузить изображение для гейта
// @Description  Загружает изображение гейта и обновляет URL в базе данных.
// @Tags         Gates
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path int true "ID гейта"
// @Param        file formData file true "Изображение гейта"
// @Success      201 {object} DTO_Resp_UploadImg
// @Failure      400 {object} map[string]string "Некорректный ID или отсутствует файл"
// @Failure      500 {object} map[string]string "Ошибка при загрузке изображения"
// @Router       /api/gates/{id}/image [post]
func (h *Handler) ApiUploadGatesImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	// Upload to object storage and update DB atomically at repo level
	imageURL, err := h.Repository.SaveServiceImage(ctx, uint(id), fileHeader)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, DTO_Resp_UploadImg{ID: id, Image: imageURL})
}

// Вспомогательная функция (не обработчик, не требует Swagger)
func generateSafeImageName(fh *multipart.FileHeader) string {
	base := fh.Filename
	ext := filepath.Ext(base)
	name := base
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	slug := re.ReplaceAllString(name, "-")
	if slug == "" {
		slug = "image"
	}
	if ext == "" {
		ext = ".png"
	}
	return slug + ext
}