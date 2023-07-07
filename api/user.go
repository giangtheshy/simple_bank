package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/giangtheshy/simple_bank/db/sqlc"
	"github.com/giangtheshy/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type loginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginResponse struct {
	SessionID uuid.UUID `json:"session_id"`
	AccessToken string       `json:"access_token"`
	AccessTokenExpireAt time.Time  `json:"access_token_expire_at"`
	RefreshToken string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time  `json:"refresh_token_expire_at"`
	User        userResponse `json:"user"`
}
type userResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	fmt.Println(req.Password)
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:     req.Username,
		FullName:     req.FullName,
		HashPassword: hashPassword,
		Email:        req.Email,
	}
	// go testWorker(arg)
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := util.CheckPasswordHash(req.Password, user.HashPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, accessPayload,err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload,err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	session,err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		ExpiresAt: refreshPayload.ExpiredAt,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := loginResponse{
SessionID: session.ID,
		AccessToken: accessToken,
		AccessTokenExpireAt: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpireAt: refreshPayload.ExpiredAt,
		User: userResponse{
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) getUser(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
	ctx.JSON(http.StatusOK, res)
}

type listUserRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}
	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := util.Map(users, func(item db.User) userResponse {
		return userResponse{Username: item.Username, FullName: item.FullName, Email: item.Email}
	})
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("username")
	err := server.store.DeleteUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
