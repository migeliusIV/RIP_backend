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

	search := ctx.Query("query")
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

	/*
		draftFrax, _ := h.Repository.GetOrCreateDraftFrax(hardcodedUserID)
		var fraxID uint = 0
		var factorsCount int = 0

		if draftFrax != nil {
			// Загружаем связанные факторы, чтобы посчитать их количество
			fullFrax, err := h.Repository.GetFraxWithFactors(draftFrax.ID)
			if err == nil {
				fraxID = fullFrax.ID
				factorsCount = len(fullFrax.FactorsLink)
			}
		}
	*/
	ctx.HTML(http.StatusOK, "gates_list.html", gin.H{
		"gates": gates,
		"query": search,
		/*
			"fraxID":        fraxID,
			"factorsCount":  factorsCount,
		*/
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

	ctx.HTML(http.StatusOK, "properties.html", gate)
}
