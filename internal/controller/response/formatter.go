package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	internalErrors "github.com/rshby/go-event-ticketing/internal/errors"
)

type ResponseDTO struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewResponse() *ResponseDTO {
	return &ResponseDTO{}
}

func (r *ResponseDTO) WithMessage(msg string) *ResponseDTO {
	r.Message = msg
	return r
}

func (r *ResponseDTO) WithData(data any) *ResponseDTO {
	r.Data = data
	return r
}

func (r *ResponseDTO) ToResponseAPI(c *gin.Context, code int) {
	c.JSON(code, r)
}

var errResponse = map[error]*internalErrors.InternalError{
	internalErrors.ErrBadRequest:    internalErrors.NewInternalError().WithCode(http.StatusBadRequest).WithMessage(internalErrors.ErrBadRequest.Error()),
	internalErrors.ErrEventNotFound: internalErrors.NewInternalError().WithCode(http.StatusNotFound).WithMessage(internalErrors.ErrEventNotFound.Error()),
}

func ResponseError(c *gin.Context, err error) {
	// check err from response error
	customErr, ok := errResponse[err]
	if !ok {
		customErr = internalErrors.NewInternalError().WithCode(http.StatusInternalServerError).WithMessage(internalErrors.ErrInternalServer.Error())
	}

	c.AbortWithStatusJSON(customErr.Code(), &ResponseDTO{
		Message: customErr.Error(),
		Data:    nil,
	})
}
