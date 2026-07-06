package interceptor

import (
	"context"
	"errors"

	errorv1 "github.com/shah-dhwanil/grpc-chat/packages/api/gen/error/v1"
	"github.com/shah-dhwanil/grpc-chat/packages/pkgerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)


func ErrorInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
)(any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		var appError *pkgerror.AppError
		if ok := errors.As(err, &appError); ok {
			var st *status.Status
			switch appError.Type {
			case pkgerror.Validation:
				st = status.New(codes.InvalidArgument, appError.Detail)
			case pkgerror.ResourceNotFound:
				st = status.New(codes.NotFound, appError.Detail)
			case pkgerror.ResourceAlreadyExists:
				st = status.New(codes.AlreadyExists, appError.Detail)
			case pkgerror.Unauthorized:
				st = status.New(codes.Unauthenticated, appError.Detail)
			case pkgerror.Internal:
				st = status.New(codes.Internal, appError.Detail)
			case pkgerror.Unknown:
				st = status.New(codes.Unknown, appError.Detail)
			default:
				st = status.New(500, appError.Detail)
			}
			var ct *structpb.Struct = nil
			if appError.Type != pkgerror.Internal && appError.Type != pkgerror.Unknown {
				ct, err = structpb.NewStruct(appError.Context)
				if err != nil {
				    return nil,err
				}
			}
			st,err := st.WithDetails(&errorv1.Error{
				Type: string(appError.Type),
				Name: appError.Name,
				Detail: appError.Detail,
				Context: ct,
			})
			if err != nil {
                return nil, status.Error(codes.Internal, "failed to marshal error")
            }
			return nil,st.Err()
		}
		
	}
	return resp, nil
}