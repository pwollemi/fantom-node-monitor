package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flashguru-git/node-monitor/log"
)

var (
	signals chan os.Signal
)

type Server struct {
}

func (s *Server) StartMonitor() {
	monitor := NewMonitor(s, time.Second*time.Duration(cfg.GetInt("monitoring_cycle")), time.Now().UTC(), time.Time{})
	go monitor.Run()
}

func NewServer() *Server {
	return &Server{}
}

func Start() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	server := NewServer()
	server.StartMonitor()

	signalReceived := <-signals
	log.WithFields(log.Fields{
		"EventName": "server_stopped",
		"Signal":    signalReceived,
	}).Infof("Shutting down server: %s", signalReceived)
}
