package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/leonkaihao/msggw/pkg/config"
	"github.com/leonkaihao/msggw/pkg/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("msggw starts.")
	args := os.Args
	if len(args) != 2 {
		log.Error("Need the argument of config json file.")
		return
	}
	cfg, err := config.NewJsonConfigLoader().Load(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.SetLevel(logLevel(cfg.LogLevel))
	svc, err := service.NewService(cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	svc.Start()
	defer svc.Close()

	excepSig := make(chan os.Signal, 1)
	signal.Notify(excepSig, os.Interrupt, syscall.SIGTERM)
	<-excepSig

	log.Info("msggw quits.")
}

func logLevel(level string) log.Level {
	l := strings.ToLower(level)
	switch l {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	default:
		return log.InfoLevel
	}
}
