package res

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DTO interface{}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationBindResponse struct {
	Status  int               `json:"status"`
	Ok      bool              `json:"ok"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

// ValidateRequest version two bind and validate by generic and custom error message
func ValidateRequest[T DTO](c echo.Context, input *T) (*ValidationBindResponse, error) {
	inputValue := reflect.ValueOf(input).Elem()
	if inputValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}

	if err := c.Bind(input); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %v", err)
	}

	var hasValidateTag bool
	inputType := reflect.TypeOf(input).Elem()
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if tag := field.Tag.Get("validate"); tag != "" {
			hasValidateTag = true
			break
		}
	}
	if !hasValidateTag {
		return nil, fmt.Errorf("این struct تگ validate ندارد و نمی‌تواند به عنوان DTO استفاده شود")
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		var errors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getValidateErrMsg(err),
			})
		}

		return &ValidationBindResponse{
			Status:  http.StatusBadRequest,
			Ok:      false,
			Message: "اعتبارسنجی ناموفق",
			Errors:  errors,
		}, nil
	}

	return nil, nil
}

func getValidateErrMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "این فیلد الزامی است"
	case "email":
		return "فرمت ایمیل نامعتبر است"
	case "min":
		return fmt.Sprintf("حداقل طول باید %s باشد", fe.Param())
	case "max":
		return fmt.Sprintf("حداکثر طول باید %s باشد", fe.Param())
	default:
		return fmt.Sprintf("خطا در اعتبار سنجی %s و متن خطا: %v", strings.ToLower(fe.Field()), fe.ActualTag())
	}
}

/*
	واسه فردا نزاریم عزیزان همه با هم بخونید که امشب شبه عشقه که امشب شب عشقه
	عزیزم دوستت دارم تولدت
*/
