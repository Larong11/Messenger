package errors

import (
	"errors"
	"net/http"
)

func HandleHTTPError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		// маппим код ошибки на HTTP статус
		httpStatus := map[ErrorCode]int{
			CodeBadRequest: http.StatusBadRequest,
			CodeInternal:   http.StatusInternalServerError,
		}[appErr.Code]

		http.Error(w, appErr.Message, httpStatus)
		return
	}

	// на случай неожиданных ошибок
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
