package errors

import (
	"net/http"

	common "github.com/gatepoint/gatepoint/api/common/v1"
	"google.golang.org/grpc/codes"
)

// ErrMap is a map of error type to CommonError
var ErrMap = map[common.ErrType]GatepointError{
	common.ErrType_ERR_TYPE_UNAUTHORIZED: {
		httpCode: http.StatusUnauthorized,
		grpcCode: codes.Unauthenticated,
		code:     common.ErrType_ERR_TYPE_UNAUTHORIZED,
	},
	common.ErrType_ERR_TYPE_NOT_FOUND: {
		httpCode: http.StatusNotFound,
		grpcCode: codes.NotFound,
		code:     common.ErrType_ERR_TYPE_NOT_FOUND,
	},
	common.ErrType_ERR_TYPE_INVALID_ARG: {
		httpCode: http.StatusBadRequest,
		grpcCode: codes.InvalidArgument,
		code:     common.ErrType_ERR_TYPE_INVALID_ARG,
	},
	common.ErrType_ERR_TYPE_PERMISSION: {
		httpCode: http.StatusForbidden,
		grpcCode: codes.PermissionDenied,
		code:     common.ErrType_ERR_TYPE_PERMISSION,
	},
	common.ErrType_ERR_TYPE_INTERNAL: {
		httpCode: http.StatusInternalServerError,
		grpcCode: codes.Internal,
		code:     common.ErrType_ERR_TYPE_INTERNAL,
	},
	common.ErrType_ERR_TYPE_ALREADY_EXIST: {
		httpCode: http.StatusConflict,
		grpcCode: codes.AlreadyExists,
		code:     common.ErrType_ERR_TYPE_ALREADY_EXIST,
	},
}
