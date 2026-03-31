package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_middleware "github.com/DimaKirejko/todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux // <- мультіплексер
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router), // remove prefix 6:15
		)
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("START HTTP SERVER!", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe() // block

		if !errors.Is(err, http.ErrServerClosed) { // server return ErrServerClosed in positive case :)
			ch <- err
		}

	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and servre HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server..")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("Brutal Shutdown HTTP server %w", err)
		}

		h.log.Warn("HTTP server stoped")
	}

	return nil
}
