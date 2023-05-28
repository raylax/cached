package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raylax/cached/redis"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(context.Background())

	logger := hclog.L()

	eg.Go(func() error {
		redisService := redis.NewService(ctx, logger.Named("REDIS"), ":1234")
		return redisService.Start()
	})

	eg.Go(func() error {
		s := <-signalCh
		if s != nil {
			return fmt.Errorf("received signal [%s]", s.String())
		}
		return nil
	})

	err := eg.Wait()
	if err != nil {
		logger.Info("exiting", "reason", err.Error())
	}

}
