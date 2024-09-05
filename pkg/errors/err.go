package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/gatepoint/gatepoint/api/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatepointError struct {
	code     common.ErrType
	param    []interface{}
	httpCode int
	grpcCode codes.Code

	//time    string
	//message string
	//detail  string
	error error
}

func (e GatepointError) HTTPCode() int {
	if e.httpCode != 0 {
		return e.httpCode
	}
	return http.StatusInternalServerError
}

func (e GatepointError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&common.Error{
		Code:    e.code,
		Message: e.Error(),
	})
}

func (e GatepointError) GRPCStatus() *status.Status {
	s := status.New(e.grpcCode, e.Error())
	s, _ = s.WithDetails(&common.Error{
		Code:    e.code,
		Message: e.Error(),
		Detail:  fmt.Sprintf("%v", e.param),
	})
	return s
}

func (e GatepointError) Error() string {
	return e.error.Error()
}

func (e GatepointError) WithError(err error) GatepointError {
	e.error = err
	return e
}

func (e GatepointError) Params(params ...interface{}) GatepointError {
	e.param = params
	return e
}

func ToGatepointError(e error) GatepointError {
	var g GatepointError
	var params []interface{}

	if s, ok := status.FromError(e); ok {
		g = ErrMap[common.ErrType_ERR_TYPE_INTERNAL].WithError(e).Params(e.Error())
		if details := s.Details(); len(details) > 0 {
			if v, ok := details[0].(*common.Error); ok {
				g = ErrMap[v.Code].WithError(e).Params(v.GetDetail())
			}
		} else {
			params = append(params, e.Error())
		}
	} else {
		g = ErrMap[common.ErrType_ERR_TYPE_INTERNAL].WithError(e).Params(e)
	}
	return g
}
