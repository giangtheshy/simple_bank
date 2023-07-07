package api

import (
	"github.com/gin-gonic/gin"
)

const ApiPrefix = "/api/v1"

func (server *Server) SetupRoutes() {
	routerDefault := gin.Default()
	routerDefault.SetTrustedProxies(nil)

	v1 := routerDefault.Group(ApiPrefix)

	router := v1.Group("/")
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Auth
	{
		router.POST("/users/login", server.loginUser)
		authRoutes.POST("/users/", server.createUser)
		router.POST("/users/renew_token", server.renewAccessToken)
		authRoutes.GET("/users/:username", server.getUser)
		authRoutes.GET("/users/", server.getUsers)
	}

	// Account
	{

		authRoutes.POST("/accounts/", server.createAccount)
		authRoutes.GET("/accounts/:id", server.getAccount)
		router.GET("/accounts/", server.getAccounts)
	}

	// Transfer
	{
		authRoutes.POST("/transfer/", server.createTransfer)
	}
	server.router=routerDefault
}
