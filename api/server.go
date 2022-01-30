package api

import (
	db "github.com/ferueda/simplebank-go/db/sqlc"
	"github.com/ferueda/simplebank-go/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store *db.Store, tm token.Maker) (*Server, error) {
	s := Server{store: store, tokenMaker: tm}
	r := gin.Default()

	r.POST("/accounts", s.createAccount)
	r.GET("/accounts", s.listAccounts)
	r.GET("/accounts/:id", s.getAccount)
	r.DELETE("/accounts/:id", s.deleteAccount)

	r.POST("/transfers", s.createTransfer)
	r.GET("/transfers", s.listTransfers)
	r.GET("/transfers/:id", s.getTransfer)

	r.POST("/users", s.createUser)
	r.POST("/users/login", s.loginUser)

	s.router = r
	return &s, nil
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
