package errors

import (
	errorHelper "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	includeErrorDetails = true
)

func TestTranslateToHttpError404_nil_nil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serviceError := NewServiceErrorNotFound(nil, "")

	errStr := eris.ToString(serviceError, true)
	fmt.Println(errStr)
	/*
		SERVICE_ERROR_NOT_FOUND
			errors.TestTranslateToHttpError404_nil_nil:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/rest_api_error_test.go:23
			errors.NewServiceErrorNotFound:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:142
			errors.newServiceError:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:121
		SERVICE_ERROR_NOT_FOUND
	*/

	TranslateServiceErrorToAPIError(c, serviceError, includeErrorDetails)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTranslateToHttpError404_x_nil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serviceError := NewServiceErrorNotFound(errorHelper.New("fromBackend404"), "")

	errStr := eris.ToString(serviceError, true)
	fmt.Println(errStr)

	/*
		SERVICE_ERROR_NOT_FOUND
			errors.TestTranslateToHttpError404_x_nil:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/rest_api_error_test.go:44
			errors.NewServiceErrorNotFound:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:142
			errors.newServiceError:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:125
		fromBackend404
	*/

	TranslateServiceErrorToAPIError(c, serviceError, includeErrorDetails)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTranslateToHttpError404_nil_x(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serviceError := NewServiceErrorNotFound(nil, "custom404")

	errStr := eris.ToString(serviceError, true)
	fmt.Println(errStr)
	/*
		custom404
			errors.TestTranslateToHttpError404_nil_x:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/rest_api_error_test.go:66
			errors.NewServiceErrorNotFound:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:142
			errors.newServiceError:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:129
		SERVICE_ERROR_NOT_FOUND
	*/

	TranslateServiceErrorToAPIError(c, serviceError, includeErrorDetails)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTranslateToHttpError404_x_x(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serviceError := NewServiceErrorNotFound(errorHelper.New("fromBackend404"), "custom404")

	errStr := eris.ToString(serviceError, true)
	fmt.Println(errStr)
	/*
		custom404
			errors.TestTranslateToHttpError404_x_x:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/rest_api_error_test.go:87
			errors.NewServiceErrorNotFound:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:142
			errors.newServiceError:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:133
		fromBackend404
	*/

	TranslateServiceErrorToAPIError(c, serviceError, includeErrorDetails)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTranslateToHttpErrorNormalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serviceError := errorHelper.New("something bad happened")

	errStr := eris.ToString(serviceError, true)
	fmt.Println(errStr)

	/*
		something bad happened
	*/

	TranslateServiceErrorToAPIError(c, serviceError, includeErrorDetails)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
