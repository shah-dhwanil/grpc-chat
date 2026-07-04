package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	
)(any,error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	//Temp Impl
	// TODO: Replace with actual authentication logic after authentication is implemented
	ctxWithUid := context.WithValue(ctx, "user_id", md["user_id"][0])
	return handler(ctxWithUid, req)
}