package handler

import (
    "net/http"
    "front_start/internal/app/ds"

    "github.com/gin-gonic/gin"
)

// simple singleton for demo auth
var creatorUserID uint = 1

func currentUserID() uint { return creatorUserID }

// moved to serialization.go

// moved to serialization.go

func (h *Handler) ApiRegister(ctx *gin.Context) {
    var req DTO_registerRequest
    if err := ctx.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    if err := h.Repository.RegisterUser(req.Login, req.Password); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    h.okJSON(ctx, http.StatusCreated, DTO_RegisterResponse{Login: req.Login})
}

func (h *Handler) ApiMe(ctx *gin.Context) {
    user, err := h.Repository.GetUserByID(currentUserID())
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, ds.ToUserPublic(user))
}

func (h *Handler) ApiUpdateMe(ctx *gin.Context) {
    var req DTO_updateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    updated, err := h.Repository.UpdateUser(currentUserID(), req.Password)
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, ds.ToUserPublic(updated))
}

func (h *Handler) ApiLogin(ctx *gin.Context) {
    var req DTO_registerRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    // For lab: pretend login success if exists
    if err := h.Repository.CheckUser(req.Login, req.Password); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    h.okJSON(ctx, http.StatusOK, DTO_LoginResponse{Login: req.Login})
}

func (h *Handler) ApiLogout(ctx *gin.Context) {
    h.okJSON(ctx, http.StatusOK, DTO_LogoutResponse{Logout: true})
}


