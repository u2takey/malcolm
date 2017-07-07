package server

import (
	"github.com/gin-gonic/gin"

	"github.com/u2takey/malcolm/model"
)

type Service struct {
	config model.Config
}

func New(cfg *model.Config) *Service {
	svc := &Service{
		config: cfg,
	}
	return svc
}

func (s *Service) GetPing(c *gin.Context) {
	c.String(200, "pong")
}
