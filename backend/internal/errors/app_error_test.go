package errors

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type validationFixture struct {
	Required  string   `json:"required" validate:"required"`
	Email     string   `json:"email" validate:"email"`
	Min       string   `json:"min" validate:"min=3"`
	Max       string   `json:"max" validate:"max=5"`
	Len       string   `json:"len" validate:"len=4"`
	Numeric   string   `json:"numeric" validate:"numeric"`
	Alpha     string   `json:"alpha" validate:"alpha"`
	Alphanum  string   `json:"alphanum" validate:"alphanum"`
	URL       string   `json:"url" validate:"url"`
	UUID      string   `json:"uuid" validate:"uuid"`
	OneOf     string   `json:"oneof" validate:"oneof=red green blue"`
	Gte       int      `json:"gte" validate:"gte=10"`
	Lte       int      `json:"lte" validate:"lte=5"`
	Gt        int      `json:"gt" validate:"gt=10"`
	Lt        int      `json:"lt" validate:"lt=5"`
	Eq        int      `json:"eq" validate:"eq=7"`
	Ne        int      `json:"ne" validate:"ne=7"`
	Unique    []string `json:"unique" validate:"unique"`
	OmitEmpty string   `json:"omit_empty" validate:"omitempty"`
	Custom    string   `json:"custom" validate:"customrule"`
}

func TestNewAppError(t *testing.T) {
	err := NewAppError("CODE", "message", ErrorTypeValidation, http.StatusBadRequest)
	require.Equal(t, "CODE", err.Code)
	require.Equal(t, "message", err.Message)
	require.Equal(t, ErrorTypeValidation, err.Type)
	require.Equal(t, http.StatusBadRequest, err.HTTPStatus)
	require.NotZero(t, err.Timestamp)
	require.NotEmpty(t, err.StackTrace)
}

func TestAppError_Error(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := NewAppError(constants.ValidationError, "invalid input", ErrorTypeValidation, http.StatusBadRequest)
		require.Equal(t, "VALIDATION_ERROR: invalid input", err.Error())
	})

	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("underlying")
		err := WrapError(cause, constants.ValidationError, "invalid input", ErrorTypeValidation, http.StatusBadRequest)
		require.Contains(t, err.Error(), "VALIDATION_ERROR: invalid input")
		require.Contains(t, err.Error(), "underlying")
	})
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := WrapError(cause, constants.InternalError, "wrapped", ErrorTypeInternal, http.StatusInternalServerError)
	require.Equal(t, cause, err.Unwrap())
}

func TestAppError_WithContextOperationResource(t *testing.T) {
	err := NewAppError(constants.NotFound, "missing", ErrorTypeNotFound, http.StatusNotFound)
	err = err.WithContext("id", "123").
		WithOperation("get_user").
		WithResource("user")

	require.Equal(t, "123", err.Context["id"])
	require.Equal(t, "get_user", err.Operation)
	require.Equal(t, "user", err.Resource)

	err2 := NewAppError(constants.NotFound, "missing", ErrorTypeNotFound, http.StatusNotFound)
	err2 = err2.WithContext("key", "value")
	require.NotNil(t, err2.Context)
}

func TestErrorConstructors(t *testing.T) {
	cause := errors.New("cause")

	tests := []struct {
		name       string
		err        *AppError
		code       string
		errType    ErrorType
		httpStatus int
	}{
		{
			name:       "ValidationError",
			err:        ValidationError("validation failed", cause),
			code:       constants.ValidationError,
			errType:    ErrorTypeValidation,
			httpStatus: http.StatusBadRequest,
		},
		{
			name:       "NotFoundError",
			err:        NotFoundError("user", cause),
			code:       constants.NotFound,
			errType:    ErrorTypeNotFound,
			httpStatus: http.StatusNotFound,
		},
		{
			name:       "UnauthorizedError",
			err:        UnauthorizedError("denied", cause),
			code:       constants.Unauthorized,
			errType:    ErrorTypeUnauthorized,
			httpStatus: http.StatusUnauthorized,
		},
		{
			name:       "CodedUnauthorized",
			err:        CodedUnauthorized("CUSTOM_UNAUTHORIZED", cause),
			code:       "CUSTOM_UNAUTHORIZED",
			errType:    ErrorTypeUnauthorized,
			httpStatus: http.StatusUnauthorized,
		},
		{
			name:       "CodedValidation",
			err:        CodedValidation("CUSTOM_VALIDATION", cause),
			code:       "CUSTOM_VALIDATION",
			errType:    ErrorTypeValidation,
			httpStatus: http.StatusBadRequest,
		},
		{
			name:       "CodedConflict",
			err:        CodedConflict("CUSTOM_CONFLICT", cause),
			code:       "CUSTOM_CONFLICT",
			errType:    ErrorTypeConflict,
			httpStatus: http.StatusConflict,
		},
		{
			name:       "ForbiddenError",
			err:        ForbiddenError("forbidden", cause),
			code:       constants.Forbidden,
			errType:    ErrorTypeForbidden,
			httpStatus: http.StatusForbidden,
		},
		{
			name:       "ConflictError",
			err:        ConflictError("conflict", cause),
			code:       constants.BadRequest,
			errType:    ErrorTypeConflict,
			httpStatus: http.StatusConflict,
		},
		{
			name:       "InternalError",
			err:        InternalError("internal", cause),
			code:       constants.InternalError,
			errType:    ErrorTypeInternal,
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:       "DatabaseError",
			err:        DatabaseError("db failed", cause),
			code:       constants.DatabaseError,
			errType:    ErrorTypeDatabase,
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:       "ExternalServiceError",
			err:        ExternalServiceError("upstream failed", cause),
			code:       constants.ExternalServiceError,
			errType:    ErrorTypeExternal,
			httpStatus: http.StatusBadGateway,
		},
		{
			name:       "CacheError",
			err:        CacheError("cache failed", cause),
			code:       constants.InternalError,
			errType:    ErrorTypeCache,
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:       "TimeoutError",
			err:        TimeoutError("timed out", cause),
			code:       constants.InternalError,
			errType:    ErrorTypeTimeout,
			httpStatus: http.StatusRequestTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.code, tt.err.Code)
			require.Equal(t, tt.errType, tt.err.Type)
			require.Equal(t, tt.httpStatus, tt.err.HTTPStatus)
			require.Equal(t, cause, tt.err.Cause)
			require.NotEmpty(t, tt.err.StackTrace)
		})
	}
}

func TestValidationErrorWithDetails(t *testing.T) {
	cause := errors.New("validation")
	err := ValidationErrorWithDetails("invalid", cause, map[string]string{
		"email": "must be valid",
		"name":  "required",
	})

	require.Equal(t, constants.ValidationError, err.Code)
	require.Equal(t, "must be valid", err.Context["email"])
	require.Equal(t, "required", err.Context["name"])
}

func TestParseValidationErrors(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.RegisterValidation("customrule", func(_ validator.FieldLevel) bool {
		return false
	}))

	input := validationFixture{
		Email:    "not-an-email",
		Min:      "ab",
		Max:      "toolong",
		Len:      "x",
		Numeric:  "abc",
		Alpha:    "123",
		Alphanum: "!!!",
		URL:      "not-url",
		UUID:     "not-uuid",
		OneOf:    "yellow",
		Gte:      1,
		Lte:      10,
		Gt:       5,
		Lt:       10,
		Eq:       1,
		Ne:       7,
		Unique:   []string{"a", "a"},
		Custom:   "bad",
	}

	validateErr := v.Struct(input)
	require.Error(t, validateErr)

	fieldErrors := ParseValidationErrors(validateErr)
	require.Equal(t, "Required is required", fieldErrors["Required"])
	require.Equal(t, "Email must be a valid email address", fieldErrors["Email"])
	require.Equal(t, "Min must be at least 3 characters long", fieldErrors["Min"])
	require.Equal(t, "Max must be at most 5 characters long", fieldErrors["Max"])
	require.Equal(t, "Len must be exactly 4 characters long", fieldErrors["Len"])
	require.Equal(t, "Numeric must be a valid number", fieldErrors["Numeric"])
	require.Equal(t, "Alpha must contain only letters", fieldErrors["Alpha"])
	require.Equal(t, "Alphanum must contain only letters and numbers", fieldErrors["Alphanum"])
	require.Equal(t, "URL must be a valid URL", fieldErrors["URL"])
	require.Equal(t, "UUID must be a valid UUID", fieldErrors["UUID"])
	require.Equal(t, "OneOf must be one of: red green blue", fieldErrors["OneOf"])
	require.Equal(t, "Gte must be greater than or equal to 10", fieldErrors["Gte"])
	require.Equal(t, "Lte must be less than or equal to 5", fieldErrors["Lte"])
	require.Equal(t, "Gt must be greater than 10", fieldErrors["Gt"])
	require.Equal(t, "Lt must be less than 5", fieldErrors["Lt"])
	require.Equal(t, "Eq must be equal to 7", fieldErrors["Eq"])
	require.Equal(t, "Ne must not be equal to 7", fieldErrors["Ne"])
	require.Equal(t, "Unique must be unique", fieldErrors["Unique"])
	require.Contains(t, fieldErrors["Custom"], "Custom is invalid")

	omitemptyErr := validator.ValidationErrors{
		stubFieldError{field: "SkipMe", tag: "omitempty", value: ""},
	}
	require.Empty(t, ParseValidationErrors(omitemptyErr))
}

func TestParseValidationErrors_NonValidationError(t *testing.T) {
	require.Empty(t, ParseValidationErrors(errors.New("plain error")))
}

func TestAppErrorHelpers(t *testing.T) {
	appErr := ValidationError("bad request", nil)
	require.True(t, IsAppError(appErr))
	require.Equal(t, appErr, GetAppError(appErr))
	require.Equal(t, constants.ValidationError, GetErrorCode(appErr))
	require.Equal(t, http.StatusBadRequest, GetHTTPStatus(appErr))
	require.Equal(t, "bad request", GetErrorMessage(appErr))

	plain := errors.New("plain failure")
	require.False(t, IsAppError(plain))
	require.Nil(t, GetAppError(plain))
	require.Equal(t, constants.InternalError, GetErrorCode(plain))
	require.Equal(t, http.StatusInternalServerError, GetHTTPStatus(plain))
	require.Equal(t, "plain failure", GetErrorMessage(plain))
}

type stubFieldError struct {
	field, tag, param string
	value             interface{}
}

func (s stubFieldError) Tag() string                      { return s.tag }
func (s stubFieldError) ActualTag() string                { return s.tag }
func (s stubFieldError) Namespace() string                { return "" }
func (s stubFieldError) StructNamespace() string          { return "" }
func (s stubFieldError) Field() string                    { return s.field }
func (s stubFieldError) StructField() string              { return s.field }
func (s stubFieldError) Value() interface{}               { return s.value }
func (s stubFieldError) Param() string                    { return s.param }
func (s stubFieldError) Kind() reflect.Kind               { return reflect.String }
func (s stubFieldError) Type() reflect.Type               { return reflect.TypeOf("") }
func (s stubFieldError) Translate(_ ut.Translator) string { return "" }
func (s stubFieldError) Error() string                    { return s.field + " failed" }
