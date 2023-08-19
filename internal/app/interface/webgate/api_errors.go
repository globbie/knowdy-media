package webgate

import (
	"github.com/globbie/knowdy-media/internal/app/usecases/upload"
	"net/http"
)

type HTTPStatus int
type errorCodeMap map[error]HTTPStatus

var errorCodes = errorCodeMap{
	upload.ErrUnauthorized: http.StatusUnauthorized,
	upload.ErrAlreadyExists: http.StatusConflict,
	upload.ErrRepoNotFound: http.StatusNotFound,
}

func (wg *WebGate) ConvertToHTTPStatus(err error) int {
	httpStatus, ok := errorCodes[err]
	if !ok {
		return http.StatusInternalServerError
	}
	return int(httpStatus)
}

