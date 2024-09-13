package tcp

import (
	"context"
	"fmt"
	"github.com/atomwqh/MyGodis/interface/tcp"
	"github.com/atomwqh/MyGodis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/*
a tcp server
*/

type Config struct {
	Address string        `yaml:"address"`
	MaxConn int           `yaml:"max_conn"`
	Timeout time.Duration `yaml:"timeout"`
}

// ClientCounter record current number of clients
var ClientCounter int32

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening ...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan chan struct{}) {
	errCh := make(chan error, 1)
	defer close(errCh)
	go func() {
		select {
		case <-closeChan:
			logger.Info("get close signal")
		case err := <-errCh:
			logger.Error(fmt.Sprintf("accept error: %s", err.Error()))
		}
		logger.Info("accept close")
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				logger.Infof("accept occours temporary error: %v, retry in 5ms", err)
				time.Sleep(time.Millisecond * 5)
				continue
			}
			errCh <- err
			break
		}
		// handle
		logger.Info("accept link")
		ClientCounter++
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
				atomic.AddInt32(&ClientCounter, -1)
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
