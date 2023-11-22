package main

import (
	"cian-parse/internals/app"
	"cian-parse/internals/config"
	"cian-parse/pkg/logger"
	"context"
	"os"
	"os/signal"
)

func main() {
	cfg := config.SetupConfig()
	log := logger.SetupLogger()
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(ctx, cfg, log)

	go func() {
		oscall := <-c
		log.Info("system call:%+v", oscall)
		server.ShotDown()
		cancel()
	}()

	server.Serve()

	/*p := parse.NewParser(log)
	err = p.ParseImmovable()
	if err != nil {
		log.Error("failed to parse all data")
		return
	}*/
}
