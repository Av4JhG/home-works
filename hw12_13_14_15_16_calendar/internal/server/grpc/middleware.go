package grpc

import (
	"context"
	"time"

	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryLoggingInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		clientIP := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}
		userAgent := "unknown"
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}

		logger.Infof("Received request: %s %s %s  at %s",
			info.FullMethod,
			clientIP,
			userAgent,
			time.Now().Format(time.RFC3339),
		)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		logger.Infof("Processed request: %s, latency: %s", info.FullMethod, duration)
		return resp, err
	}
}
