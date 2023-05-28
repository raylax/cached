package redis

import (
	"context"
	"github.com/hashicorp/go-hclog"
)

type Service struct {
	addr    string
	context context.Context
	s       *server
	logger  hclog.Logger
}

func NewService(context context.Context, logger hclog.Logger, addr string) *Service {
	return &Service{
		addr:    addr,
		context: context,
		logger:  logger,
	}
}

func (s *Service) Start() error {
	s.s = newServer(s.logger, s.addr)
	go func() {
		<-s.context.Done()
		_ = s.Stop()
	}()
	s.logger.Info("starting redis server", "server", s.addr)
	return s.s.start()
}

func (s *Service) Stop() error {
	s.logger.Info("stopping redis server", "addr", s.addr)
	return s.s.stop()
}
