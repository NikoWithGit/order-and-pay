package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type server struct {
	Gin *gin.Engine
}

func NewServer() *server {
	return &server{gin.Default()}
}

func (s *server) Start() error {
	if s.Gin == nil {
		return errors.New("ERROR:No Gin Engine")
	}
	return s.Gin.Run(":8080")
}
