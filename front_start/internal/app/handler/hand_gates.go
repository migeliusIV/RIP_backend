package handler

import (
	"front_start/internal/app/ds"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

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
			ID_gate: gate.ID_gate,
			Title: gate.Title,
			Description: gate.Description,
			Status: gate.Status,
			Image: gate.Image,
			I0j0: gate.I0j0,
			I0j1: gate.I0j1,
			I1j0: gate.I1j0,
			I1j1: gate.I1j1,
			Matrix_koeff: gate.Matrix_koeff,
			// subject area
			FullInfo: gate.FullInfo,
			TheAxis: gate.TheAxis,
		})
	}
    ctx.JSON(http.StatusOK, represent_gates)
}

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
		Title:       req.Title,
		Description: req.Description,
		FullInfo:    req.FullInfo,
		TheAxis:     req.TheAxis,
		I0j0: req.I0j0,
		I0j1: req.I0j1,
		I1j0: req.I1j0,
		I1j1: req.I1j1,
		Matrix_koeff: req.Matrix_koeff, 
	}
	if req.Status != nil {
		gate.Status = *req.Status
	}
	if err := h.Repository.AddGate(&gate); err != nil { //CreateService -> AddGate
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, gate)
}

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
	if err := h.Repository.AddGateToTask(userID, uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
    ctx.JSON(http.StatusCreated, DTO_Resp_TaskServiceLink{TaskID: task.ID_task, ServiceID: id})
}

func (h *Handler) ApiGetCurrQTask(ctx *gin.Context) {
	draftTask, _ := h.Repository.GetDraftTask(hardcodedUserID)
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

// ApiUploadServiceImage: accepts multipart form with field "file", sets image name
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

func generateSafeImageName(fh *multipart.FileHeader) string {
	base := fh.Filename
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
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
