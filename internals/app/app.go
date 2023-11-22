package app

import (
	"cian-parse/api"
	"cian-parse/internals/app/db"
	"cian-parse/internals/app/handlers"
	"cian-parse/internals/app/processors"
	"cian-parse/internals/config"
	"cian-parse/pkg/client/mongodb"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
)

type Server struct {
	srv  *http.Server
	cfg  *config.Config
	log  *slog.Logger
	mgdb *mongo.Database
	ctx  context.Context
}

func NewServer(ctx context.Context, cfg *config.Config, log *slog.Logger) *Server {
	server := new(Server)
	server.cfg = cfg
	server.log = log
	server.ctx = ctx
	return server
}

func (s *Server) Serve() {
	s.log.Info("starting server")
	mongoDBClient, err := mongodb.NewClient(s.ctx, s.cfg.Storage.Host, s.cfg.Storage.Port, s.cfg.Storage.Username,
		s.cfg.Storage.Password, s.cfg.Storage.Database, s.cfg.Storage.AuthDB)
	if err != nil {
		panic(err)
	}
	s.mgdb = mongoDBClient
	s.log.Info("connected to mongodb")

	immovablesStorage := db.NewImmovablesStorage(s.mgdb, s.cfg.Storage.Collection, s.log)
	immovablesProcessor := processors.NewImmovablesProcessor(immovablesStorage)
	immovablesHandler := handlers.NewImmovablesHandler(s.ctx, immovablesProcessor)

	routes := api.CreateRoutes(immovablesHandler)

	/*finded, err := immovablesProcessor.FindAll(s.ctx)
	fmt.Println(finded)

	var immovable1 = models.Immovable{
		Title:          "test",
		Link:           "test",
		Data:           "test",
		Price:          4,
		PriceInitially: 2,
	}
	immovable1ID, err := immovablesProcessor.Create(context.Background(), immovable1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(immovable1ID)

	var immovable2 = models.Immovable{
		Title: "fun",
	}

	err = immovablesProcessor.Update(context.Background(), immovable1ID, immovable2)
	if err != nil {
		panic(err)
	}
	fmt.Println(immovablesProcessor.FindOne(context.Background(), "6523bcadef62f7da34ab115d"))*/

	s.srv = &http.Server{
		Addr:    ":" + s.cfg.Listen.Port,
		Handler: routes,
	}

	s.log.Info("server started in", s.cfg.Storage.Host, s.cfg.Listen.Port)

	err = s.srv.ListenAndServe()
	if err != nil {
		s.log.Error("error", err)
	}

	return
}

func (s *Server) ShotDown() {
	s.log.Info("server stopped")

	ctxShutDown, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	var err error
	if err = s.srv.Shutdown(ctxShutDown); err != nil {
		s.log.Error("server Shutdown failed", err)
	}

	s.log.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
