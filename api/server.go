package api

import (
	db "github.com/ferueda/simplebank-go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	s := Server{store: store}
	r := gin.Default()

	r.POST("/accounts", s.createAccount)
	r.GET("/accounts", s.listAccounts)
	r.GET("/accounts/:id", s.getAccount)
	r.DELETE("/accounts/:id", s.deleteAccount)

	// TODO: implement entries and transfers endpoints

	s.router = r
	return &s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
