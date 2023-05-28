package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/netpoll"
	"github.com/hashicorp/go-hclog"
	"time"
)

const (
	stopTimeoutSecond = 30
)

func init() {
	_ = netpoll.DisableGopool()
}

type server struct {
	evenLoop netpoll.EventLoop
	logger   hclog.Logger
	addr     string
}

func newServer(logger hclog.Logger, addr string) *server {
	return &server{
		logger: logger,
		addr:   addr,
	}
}

func (s *server) handelRequest(ctx context.Context, conn netpoll.Connection) error {
	s.logger.Info("new request", "remote", conn.RemoteAddr(), "active", conn.IsActive())
	r := conn.Reader()
	defer func() {
		err := r.Release()
		if err != nil {
			s.logger.Error("failed to release reader", "remote", conn.RemoteAddr(), "error", err.Error())
		}
	}()
	data, err := ReadResp(r)
	if err != nil {
		s.logger.Info("failed to read request", "remote", conn.RemoteAddr(), "error", err.Error())
		return nil
	}
	b, _ := json.Marshal(data)
	s.logger.Info("request", "remote", conn.RemoteAddr(), "data", string(b))
	return nil
}

func (s *server) handlePrepare(conn netpoll.Connection) (ctx context.Context) {
	_ = conn.AddCloseCallback(s.handelClose)
	return
}

func (s *server) handelClose(conn netpoll.Connection) error {
	s.logger.Info("connection closed", "remote", conn.RemoteAddr())
	return nil
}

func (s *server) handleConnect(ctx context.Context, conn netpoll.Connection) context.Context {
	s.logger.Info("new connection", "remote", conn.RemoteAddr())
	return ctx
}

func (s *server) start() error {
	evenLoop, err := netpoll.NewEventLoop(
		s.handelRequest,
		netpoll.WithOnPrepare(s.handlePrepare),
		netpoll.WithOnConnect(s.handleConnect),
	)
	if err != nil {
		return err
	}
	s.evenLoop = evenLoop
	listener, err := netpoll.CreateListener("tcp", s.addr)
	if err != nil {
		return err
	}
	return evenLoop.Serve(listener)
}

func (s *server) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), stopTimeoutSecond*time.Second)
	defer cancel()
	err := s.evenLoop.Shutdown(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout waiting for server to stop in %d seconds", stopTimeoutSecond)
	}
	return err
}
