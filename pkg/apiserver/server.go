package apiserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/apiserver/handlers"
	"github.com/3d0c/toto-config/pkg/config"
	"github.com/3d0c/toto-config/pkg/log"
)

// APIHTTPServer structure
type APIHTTPServer struct {
	apiVersion string
	srv        *http.Server
	logger     *zap.Logger
}

// NewAPIHTTPServer API server constructor
func NewAPIHTTPServer(cfg config.Server) (*APIHTTPServer, error) {
	var (
		cert tls.Certificate
		err  error
	)

	s := &APIHTTPServer{
		apiVersion: cfg.APIVersion,
		logger:     log.TheLogger().With(zap.String("component", "APIHTTPServer")),
	}

	s.srv = &http.Server{
		Handler:      handlers.SetupRouter(cfg),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTemout,
		TLSConfig: &tls.Config{
			Certificates: make([]tls.Certificate, 0),
			NextProtos:   []string{"http/1.1"},
			MinVersion:   tls.VersionTLS12,
			//nolint: gosec
			InsecureSkipVerify: cfg.Insecure,
		},
	}

	if cfg.Address != "" {
		s.srv.Addr = cfg.Address
	}

	if cert, err = tls.LoadX509KeyPair(cfg.Certificate, cfg.PrivateKey); err != nil {
		return nil, fmt.Errorf("error loading key pair - %s", err)
	}

	s.srv.TLSConfig.Certificates = append(
		s.srv.TLSConfig.Certificates,
		cert,
	)

	return s, nil
}

// Run starts HTTP server, ctx is used for server shutdown in case if ctx is closed
func (s *APIHTTPServer) Run(ctx context.Context) {
	loggerWithField := s.logger.With(zap.String("method", "Run"))

	go func() {
		for {
			<-ctx.Done()
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
			_ = s.srv.Shutdown(shutdownCtx)
			cancelFn()
			return
		}
	}()

	loggerWithField.With(zap.String("address", s.srv.Addr), zap.String("apiVersion", s.apiVersion)).
		Info("start server")

	if err := s.srv.ListenAndServeTLS("", ""); err != nil {
		loggerWithField.Warn("http server finished with error", zap.Error(err))
	}
}
