package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/paxaf/BrandScoutTest/internal/controller"
	"github.com/paxaf/BrandScoutTest/internal/controller/middleware"
	storage "github.com/paxaf/BrandScoutTest/internal/repo/engine"
	"github.com/paxaf/BrandScoutTest/internal/usecase"
	"github.com/paxaf/BrandScoutTest/internal/worker"
)

const (
	appHost                      = "0.0.0.0"
	appPort                      = "8080"
	defaultTimeout time.Duration = 5 * time.Second
)

type App struct {
	apiServer *http.Server
	scheduler *worker.Scheduler
}

func New() (*App, error) {
	app := &App{}
	repo, err := storage.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("failed init repo: %w", err)
	}
	service := usecase.New(repo)
	handler := controller.New(service)
	http.Handle("/tasks", middleware.SimpleMiddleware(
		http.HandlerFunc(handler.GetAll),
		http.HandlerFunc(handler.Get),
		http.HandlerFunc(handler.Add)))
	http.HandleFunc("/tasks/", handler.Delete)
	addr := net.JoinHostPort(appHost, appPort)
	app.apiServer = &http.Server{
		Addr:              addr,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: defaultTimeout,
	}

	scheduler := worker.NewScheduler(service)
	app.scheduler = scheduler
	return app, nil
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("API server started successfully. " + "Address: " + app.apiServer.Addr)
		if err := app.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start server")
		}
	}()

	app.scheduler.Start()
	<-ctx.Done()
	log.Printf("Received shutdown signal")

	return nil
}

func (app *App) Close() error {
	app.scheduler.Stop()
	err := app.apiServer.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
