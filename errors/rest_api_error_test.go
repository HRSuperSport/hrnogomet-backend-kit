package errors

import (
	errorHelper "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hrsupersport/hrnogomet-backend-kit/logging"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
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

func TestErrorStackLogging(t *testing.T) {
	//gin.SetMode(gin.TestMode)
	logging.ConfigureDefaultLoggingSetup("hrnogomet-backend-kit")

	serviceError := NewServiceErrorNotFound(errorHelper.New("fromBackend404"), "custom404")

	/* contains absolute paths
	custom404
		errors.TestErrorStackLogging:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/rest_api_error_test.go:127
		errors.NewServiceErrorNotFound:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:142
		errors.newServiceError:/Users/adambezecny/dev/hrnogomet-backend-kit/errors/service_error.go:133
	fromBackend404
	*/
	fmt.Println(eris.ToString(serviceError, true))

	/* contains relative paths thanks to setting stackPathSplitter in logging.ConfigureDefaultLoggingSetup call!
	{"L":"error","stack":[{"func":"errors.newServiceError","source":"/errors/service_error.go:133"},{"func":"errors.NewServiceErrorNotFound","source":"/errors/service_error.go:142"},{"func":"errors.TestErrorStackLogging","source":"/errors/rest_api_error_test.go:127"}],"error":"custom404: fromBackend404","T":"2023-11-06T09:35:45.36983+01:00","caller":"rest_api_error_test.go:129"}
	*/
	log.Error().Stack().Err(serviceError).Msg("")
}
