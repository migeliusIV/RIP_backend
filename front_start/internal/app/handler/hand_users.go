package handler

import (
    "net/http"
    "front_start/internal/app/ds"
    "errors"
    "strings"
	"time"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// simple singleton for demo auth
var creatorUserID uint = 1

func currentUserID() uint { return creatorUserID }

func (h *Handler) Register(ctx *gin.Context) {
    var req DTO_Req_UserReg
    if err := ctx.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

    if err := h.Repository.RegisterUser(req.Login, string(hashedPassword)); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    ctx.JSON(http.StatusOK, DTO_Resp_User{Login: req.Login})
}

func (h *Handler) Login(ctx *gin.Context) {
	var req DTO_Req_UserReg
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByUsername(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		customErr := errors.New(user.Password)
		h.errorHandler(ctx, http.StatusUnauthorized, customErr)
		return
	}

	claims := ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.JWTConfig.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:      user.ID_user,
		IsModerator: user.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTConfig.Secret))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := DTO_Resp_TokenLogin{
		Token: tokenString,
		User: DTO_User{
			ID_user: user.ID_user,
			Login:  user.Login,
            IsAdmin: user.IsAdmin,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) ApiMe(ctx *gin.Context) {
    user, err := h.Repository.GetUserByID(currentUserID())
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    ctx.JSON(http.StatusOK, DTO_Resp_User{Login: user.Login})
}

func (h *Handler) ApiUpdateMe(ctx *gin.Context) {
    var req DTO_Req_UserUpd
    if err := ctx.ShouldBindJSON(&req); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }
    updated, err := h.Repository.UpdateUser(currentUserID(), req.Password)
    if err != nil {
        h.errorHandler(ctx, http.StatusInternalServerError, err)
        return
    }
    ctx.JSON(http.StatusOK, DTO_Resp_User{Login: updated.Login})
}

func (h *Handler) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("invalid header"))
		return
	}
	tokenStr := authHeader[len("Bearer "):]

	err := h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), tokenStr, h.JWTConfig.ExpiresIn)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Деавторизация прошла успешно",
	})
}

func getUserIDFromContext(ctx *gin.Context) (uint, error) {
	value, exists := ctx.Get(userCtx)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userID, ok := value.(uint)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return userID, nil
}

/*
func (h *Handler) ApiLogout(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, DTO_Resp_UserLogout{Logout: true})
}
*/