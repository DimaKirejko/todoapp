package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/DimaKirejko/todoapp/internal/core/logger"
	core_http_response "github.com/DimaKirejko/todoapp/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	xRequestIDHeader = "X-Request-ID"
)

// func RequestID() Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(requestIDHandler(next)) /// повертає middleware (функцію, яка обгортає handler)
// 	}
// }

// func requestIDHandler(next http.Handler) func(http.ResponseWriter, *http.Request) { //створює handler з логікою middleware (виконується при запиті)
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		requestID := r.Header.Get(xRequestIDHeader)
// 		if requestID == "" {
// 			requestID = uuid.NewString()
// 		}

// 		r.Header.Set(xRequestIDHeader, requestID)
// 		w.Header().Set(xRequestIDHeader, requestID)

// 		next.ServeHTTP(w, r)
// 	}
// }

func RequestID() Middleware {
	return func(next http.Handler) http.Handler { //
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(xRequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(xRequestIDHeader, requestID)
			w.Header().Set(xRequestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(xRequestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler { /// return return need recheck
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			timeBefore := time.Now()
			log.Debug(
				">>> incoming HTTP req",
				zap.String("Method", r.Method),
				zap.Time("time", timeBefore.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP req",
				zap.Int("status code", rw.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Since(timeBefore)),
			)
		})
	}
}
