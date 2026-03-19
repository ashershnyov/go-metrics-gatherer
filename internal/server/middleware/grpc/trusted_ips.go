package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CheckIP checks whether the client's ip is within the trusted subnet. If not, will return PermissionDenied.
func CheckIP(subnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if subnet == nil {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("x-real-ip")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing header: %s", "x-real-ip")
		}

		ip := net.ParseIP(values[0])
		if ip == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid client ip")
		}

		if !subnet.Contains(ip) {
			return nil, status.Errorf(codes.PermissionDenied, "ip not trusted")
		}

		return handler(ctx, req)
	}
}
