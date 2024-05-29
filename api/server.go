package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kuthumipepple/bank-backend/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	server.router = router
	return server
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}