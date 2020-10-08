package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleGrpcError(err error) error {
	switch err.(type) {
	case *InvalidEntityError, *InvalidParameterError, *MissingParameterError:
		return status.Error(codes.InvalidArgument, err.Error())
	case *InvalidContextError:
		return status.Error(codes.FailedPrecondition, err.Error())
	case *NotImplementedError:
		return status.Error(codes.Unimplemented, err.Error())
	case *ResourceExhausted:
		return status.Error(codes.ResourceExhausted, err.Error())
	case *UnauthenticatedError:
		return status.Error(codes.Unauthenticated, err.Error())
	case *UnauthorizedError:
		return status.Error(codes.PermissionDenied, err.Error())
	case *NotFoundError:
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Unavailable, err.Error())
	}
}

// Transforms an application error into a grpc error by mapping the error type to a response code
func ToGrpc(err error) error {
	return handleGrpcError(err)
}
