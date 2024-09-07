package errors

import (
	"net/http"

	commonv1 "github.com/gatepoint/gatepoint/api/common/v1"
	"google.golang.org/grpc/codes"
)

// ErrMap is a map of error type to CommonError
var ErrMap = map[commonv1.ErrType]GatepointError{
	commonv1.ErrType_ERR_TYPE_UNAUTHORIZED: {
		httpCode: http.StatusUnauthorized,
		grpcCode: codes.Unauthenticated,
		code:     commonv1.ErrType_ERR_TYPE_UNAUTHORIZED,
	},
	commonv1.ErrType_ERR_TYPE_NOT_FOUND: {
		httpCode: http.StatusNotFound,
		grpcCode: codes.NotFound,
		code:     commonv1.ErrType_ERR_TYPE_NOT_FOUND,
	},
	commonv1.ErrType_ERR_TYPE_INVALID_ARG: {
		httpCode: http.StatusBadRequest,
		grpcCode: codes.InvalidArgument,
		code:     commonv1.ErrType_ERR_TYPE_INVALID_ARG,
	},
	commonv1.ErrType_ERR_TYPE_PERMISSION: {
		httpCode: http.StatusForbidden,
		grpcCode: codes.PermissionDenied,
		code:     commonv1.ErrType_ERR_TYPE_PERMISSION,
	},
	commonv1.ErrType_ERR_TYPE_INTERNAL: {
		httpCode: http.StatusInternalServerError,
		grpcCode: codes.Internal,
		code:     commonv1.ErrType_ERR_TYPE_INTERNAL,
	},
	commonv1.ErrType_ERR_TYPE_ALREADY_EXIST: {
		httpCode: http.StatusConflict,
		grpcCode: codes.AlreadyExists,
		code:     commonv1.ErrType_ERR_TYPE_ALREADY_EXIST,
	},
}
