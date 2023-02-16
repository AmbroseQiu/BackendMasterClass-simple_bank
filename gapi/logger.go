package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	timeNow := time.Now()
	statusCode := codes.Unknown
	logger := log.Info()
	result, err := handler(ctx, req)
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	if err != nil {
		logger = log.Error().Err(err)
	}
	duration := time.Since(timeNow)
	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Dur("duration", duration).
		Int("status code", int(statusCode)).
		Str("status text", statusCode.String()).
		Msg("receive request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		timeNow := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(timeNow)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}
		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("request path", req.RequestURI).
			Dur("duration", duration).
			Int("status code", rec.StatusCode).
			Str("status text", http.StatusText(rec.StatusCode)).
			Msg("receive http request")
	})
}
