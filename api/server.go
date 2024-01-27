package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	store  db.Store
	maker  token.Maker
	config util.Config
	router *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to creating token")
	}
	server := &Server{
		store:  store,
		maker:  maker,
		config: config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setRouter()

	return server, nil
}

func (server *Server) setRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authGroups := router.Group("/").Use(authMiddleware(server.maker))

	authGroups.POST("/accounts", server.createAccount)
	authGroups.GET("/accounts/:id", server.readAccount)
	authGroups.GET("/accounts", server.listAccount)
	authGroups.PUT("/accounts", server.UpdateAccount)
	authGroups.DELETE("/accounts/:id", server.DeleteAccount)

	authGroups.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
