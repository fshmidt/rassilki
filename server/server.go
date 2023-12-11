package server

import (
	"context"
	"github.com/fshmidt/rassilki/pkg/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer      *http.Server
	RassilkaService service.Rassilka
	MessagesService service.Messages
	Context         context.Context
	Cancel          context.CancelFunc
}

func NewServer(rassilkaService service.Rassilka, messagesService service.Messages, ctx context.Context, cancel context.CancelFunc) *Server {
	return &Server{
		RassilkaService: rassilkaService,
		MessagesService: messagesService,
		Context:         ctx,
		Cancel:          cancel,
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	defer s.Cancel()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalf("ListenAndServe error: %v\n", err)
			}
		}
	}()

	go s.StartRassilkaListener(s.Context)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Shutting down server...")
	if err := s.Shutdown(s.Context); err != nil {
		logrus.Fatalf("Error during server shutdown: %v\n", err)
	}

	return s.httpServer.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) StartRassilkaListener(ctx context.Context) {
	go func() {
		var actIds, allIds, updatedIds, recreatedIds []int
		rassilkaListener := NewRassilkaListener(s.MessagesService)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				actIds, allIds, updatedIds, recreatedIds = s.checkActive()
				rassilkaListener.UpdateStatus(actIds, allIds, updatedIds, recreatedIds)
				time.Sleep(4 * time.Second)
				rassilkaListener.ResetUpdateRecreateStatus(updatedIds, recreatedIds)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (s *Server) checkActive() (actIds, allIds, updatedIds, recreatedIds []int) {

	actIds, err := s.RassilkaService.CheckActive()
	if err != nil {
		logrus.Errorf("Error getting active rassilki: %v", err)
		return nil, nil, nil, nil
	}
	allIds, err2 := s.RassilkaService.GetAll()
	if err2 != nil {
		logrus.Errorf("Error getting all rassilki id: %v", err)
		return nil, nil, nil, nil
	}
	updatedIds, err3 := s.RassilkaService.CheckUpdated()
	if err3 != nil {
		logrus.Errorf("Error getting updated rassilki id: %v", err)
		return nil, nil, nil, nil
	}
	recreatedIds, err4 := s.RassilkaService.CheckRecreated()
	if err4 != nil {
		logrus.Errorf("Error getting recreated rassilki id: %v", err)
		return nil, nil, nil, nil
	}
	return actIds, allIds, updatedIds, recreatedIds
}
