package api

import (
	"fmt"

	db "github.com/giangtheshy/simple_bank/db/sqlc"
	"github.com/giangtheshy/simple_bank/token"
	"github.com/giangtheshy/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	
	// gin.DefaultWriter = io.Discard
	// gin.DefaultErrorWriter = io.Discard
	// router.Use(Logger())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRoutes()
	return server,nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// func Logger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 			start := time.Now()

// 			c.Next()

// 			latency := time.Since(start)
// 			status := c.Writer.Status()
// 			err := c.Errors.String()
// 			if status>=400{

// 				log.Printf("[%d] \t %s \t %s \t %s \t[ERROR] \t%s", status, c.Request.Method, c.Request.URL.Path,latency,err)
// 			}else{

// 				log.Printf("[%d] \t %s \t %s \t %s ", status, c.Request.Method, c.Request.URL.Path,latency)
// 			}

// 	}
// }
