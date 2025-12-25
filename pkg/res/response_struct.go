package res

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response[T any] struct {
	OK         bool    `json:"ok"`
	StatusCode int     `json:"statusCode"`
	Message    string  `json:"message,omitempty"`
	Data       T       `json:"data,omitempty"`
	Meta       any     `json:"meta,omitempty"`
	Error      *string `json:"error,omitempty"`
}

func SuccessResponse[T any](data T, message string, statusCode ...int) *Response[T] {
	resp := &Response[T]{
		OK:      true,
		Message: message,
		Data:    data,
	}
	if len(statusCode) > 0 {
		resp.StatusCode = statusCode[0]
	} else {
		resp.StatusCode = http.StatusOK
	}
	return resp
}

func SuccessMessage(message string, statusCode ...int) *Response[struct{}] {
	return SuccessResponse(struct{}{}, message, statusCode...)
}

func ErrorResponse[T any](message string, err error, statusCode ...int) *Response[T] {
	resp := &Response[T]{
		OK:      false,
		Message: message,
	}

	if len(statusCode) > 0 {
		resp.StatusCode = statusCode[0]
	} else {
		resp.StatusCode = http.StatusInternalServerError
	}

	if err != nil {
		errMsg := "ERROR: " + err.Error()
		resp.Error = &errMsg
	}

	return resp
}

func ErrorMessage[T any](message string, statusCode ...int) *Response[T] {
	resp := &Response[T]{
		OK:      false,
		Message: message,
	}
	if len(statusCode) > 0 {
		resp.StatusCode = statusCode[0]
	} else {
		resp.StatusCode = http.StatusInternalServerError
	}
	return resp
}

func Error[T any](err error, code ...int) *Response[T] {
	errMsg := "ERROR: " + err.Error()

	resp := &Response[T]{
		OK:    false,
		Error: &errMsg,
	}

	if len(code) > 0 {
		resp.StatusCode = code[0]
	} else {
		resp.StatusCode = http.StatusInternalServerError
	}

	return resp
}

func JSON[T any](c echo.Context, r *Response[T]) error {
	return c.JSON(r.StatusCode, r)
}

// برای خطاهای اعتبارسنجی
// func ValidationErrorResponse(errors map[string]string) *Response[map[string]string] {
// 	return &Response[map[string]string]{
// 		OK:      false,
// 		Message: "Validation failed",
// 		Data:    errors,
// 	}
// }

func ValidationErrorResponse(errors []string) *Response[[]string] {
	return &Response[[]string]{
		OK:         false,
		Message:    "Validation failed",
		Data:       errors,
		StatusCode: http.StatusBadRequest,
		Error:      &errors[0],
	}
}
