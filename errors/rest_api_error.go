package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"net/http"
)

// ApiError is generic structure used to return error from REST API
type ApiError struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

// TranslateServiceErrorToAPIError holds mapping between protocol agnostic
// service errors and rest api specific errors with http status codes
func TranslateServiceErrorToAPIError(ctx *gin.Context, err error, includeDetails bool) {
	var target1 *ServiceErrorUnauthorized
	if eris.As(err, &target1) {
		ReturnUnauthorizedError(ctx, err, includeDetails)
		return
	}

	var target2 *ServiceErrorForbidden
	if eris.As(err, &target2) {
		ReturnForbiddenError(ctx, err, includeDetails)
		return
	}

	var target3 *ServiceErrorNotFound
	if eris.As(err, &target3) {
		ReturnNotFoundError(ctx, err, includeDetails)
		return
	}

	var target4 *ServiceErrorNotImplemented
	if eris.As(err, &target4) {
		ReturnNotImplementedError(ctx, err, includeDetails)
		return
	}

	ReturnInternalServerError(ctx, err, includeDetails)
}

func getGinH(err error, includeDetails bool) any {
	if includeDetails {
		detail := eris.ToString(err, true)

		return gin.H{
			"error":  err.Error(),
			"detail": detail,
		}
	} else {
		return gin.H{"error": err.Error()}
	}
}

func ReturnInternalServerError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusInternalServerError, getGinH(err, includeDetails))
}
func ReturnBadRequestError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusBadRequest, getGinH(err, includeDetails))
}
func ReturnNotImplementedError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusNotImplemented, getGinH(err, includeDetails))
}
func ReturnNotFoundError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusNotFound, getGinH(err, includeDetails))
}

func ReturnUnauthorizedError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusUnauthorized, getGinH(err, includeDetails))
}

func ReturnForbiddenError(c *gin.Context, err error, includeDetails bool) {
	c.JSON(http.StatusForbidden, getGinH(err, includeDetails))
}
