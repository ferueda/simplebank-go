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

	r.POST("/users", s.createUser)
	r.POST("/users/login", s.loginUser)

	authRoutes := r.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTransfer)
	authRoutes.GET("/transfers", s.listTransfers)
	authRoutes.GET("/transfers/:id", s.getTransfer)

	s.router = r
	return &s, nil
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
