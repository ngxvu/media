package errors

import (
	"github.com/gin-gonic/gin"
	"media-service/services/api_response"
	"net/http"
)

// Response trả về cho APP FE khi có lỗi

type MetaData struct {
	TraceID string `json:"traceId"`
	Success bool   `json:"success"`
	Status  int    `json:"status"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorDataResult struct {
	ErrorResp ErrorData `json:"error"`
}

func (e *ErrorDataResult) Error() string {
	return e.ErrorResp.Message
}

type MessagesResponse struct {
	Meta MetaData  `json:"meta"`
	Err  ErrorData `json:"error"`
}

func ErrorResponse(err string, errType string) *ErrorDataResult {
	return &ErrorDataResult{
		ErrorResp: ErrorData{
			Message: err,
			Code:    errType,
		},
	}
}

const (
	WrongPassOrUser = "account_not_found"

	// NotFound 404 error indicates a missing / not found record
	NotFound          = "NotFound"
	notFoundMessage   = "record not found"
	VnNotFoundMessage = "Không tìm thấy dữ liệu"

	// ValidationError 400 indicates an error in input validation
	ValidationError          = "ValidationError"
	ValidationErrorMessage   = "validation error"
	VnValidationErrorMessage = "Dữ liệu không hợp lệ"

	// ResourceAlreadyExists indicates a duplicate / already existing record
	ResourceAlreadyExists       = "ResourceAlreadyExists"
	alreadyExistsErrorMessage   = "resource already exists"
	VnAlreadyExistsErrorMessage = "Dữ liệu đã tồn tại"

	// RepositoryError indicates a repository (e.g database) error
	RepositoryError          = "RepositoryError"
	repositoryErrorMessage   = "error in repository operation"
	VnRepositoryErrorMessage = "Có lỗi xảy ra"

	// NotAuthenticated indicates an authentication error
	NotAuthenticated             = "NotAuthenticated"
	notAuthenticatedErrorMessage = "not Authenticated"
	VnNotAuthenticated           = "Không xác thực"

	// TokenGeneratorError indicates an token generation error
	TokenGeneratorError        = "TokenGeneratorError"
	tokenGeneratorErrorMessage = "error in token generation"
	VnTokenGeneratorMessage    = "Có lỗi xảy ra"

	// NotAuthorized indicates an authorization error
	NotAuthorized             = "NotAuthorized"
	notAuthorizedErrorMessage = "not authorized"
	VnNotAuthorizedMessage    = "Không có quyền"

	// UnknownError 500 indicates an error that the app cannot find the cause for
	UnknownError          = "UnknownError"
	UnknownErrorMessage   = "something went wrong"
	VnUnknownErrorMessage = "Có lỗi xảy ra"

	// account_not_found
	AccountNotFound          = "AccountNotFound"
	AccountNotFoundMessage   = "account not found"
	VnAccountNotFoundMessage = "Tài khoản không tồn tại"

	// missing_required_fields
	MissingRequiredFields          = "MissingRequiredFields"
	MissingRequiredFieldsMessage   = "missing required fields"
	VnMissingRequiredFieldsMessage = "Thiếu dữ liệu nhập vào "
)

// Handler is Gin middleware to handle errors.
func HandlerError(c *gin.Context) {
	// Execute request handlers and then handle any errors
	c.Next()
	errs := c.Errors

	if len(errs) > 0 {
		err, ok := errs[0].Err.(*ErrorDataResult)
		if ok {
			meta := api_response.NewMetaData()

			var status int
			switch err.ErrorResp.Code {
			case NotFound:
				status = http.StatusNotFound
			case ValidationError:
				status = http.StatusBadRequest
			case ResourceAlreadyExists:
				status = http.StatusConflict
			case NotAuthenticated:
				status = http.StatusUnauthorized
			case NotAuthorized:
				status = http.StatusForbidden
			case RepositoryError, UnknownError, TokenGeneratorError:
				status = http.StatusInternalServerError
			case AccountNotFound:
				status = http.StatusNotFound
			case MissingRequiredFields:
				status = http.StatusBadRequest
			default:
				status = http.StatusInternalServerError
			}

			resp := MessagesResponse{
				Meta: MetaData{
					TraceID: meta.TraceID,
					Success: false,
					Status:  status,
				},
				Err: ErrorData{
					Code:    err.ErrorResp.Code,
					Message: err.ErrorResp.Message,
				},
			}

			c.JSON(status, resp)
			return
		}
	}
}
