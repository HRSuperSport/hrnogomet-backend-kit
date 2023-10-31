package errors

import "github.com/rotisserie/eris"

// ServiceError represents protocol agnostic business error that might be raised within
// service or repository. This error than must be translated at handler level to protocol
// specific error, e.g. HTTP status code with response payload in case of REST API
// Default implementation for REST API translation see ApiError and TranslateServiceErrorToAPIError
type ServiceError interface {
	SetErrorText(text string)
	SetNestedError(err error)
	Error() string // ServiceError is also Error, i.e. implementors of ServiceError implement implicitly also Error interface
}

type ServiceErrorUnauthorized struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorUnauthorized) Error() string {
	return e.ErrorText
}

func (e *ServiceErrorUnauthorized) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorUnauthorized) SetNestedError(err error) {
	e.NestedError = err
}

type ServiceErrorNotFound struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorNotFound) Error() string {
	return e.ErrorText
}

func (e *ServiceErrorNotFound) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorNotFound) SetNestedError(err error) {
	e.NestedError = err
}

type ServiceErrorForbidden struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorForbidden) Error() string {
	return e.ErrorText
}

func (e *ServiceErrorForbidden) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorForbidden) SetNestedError(err error) {
	e.NestedError = err
}

type ServiceErrorBadRequest struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorBadRequest) Error() string {
	return e.ErrorText
}

func (e *ServiceErrorBadRequest) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorBadRequest) SetNestedError(err error) {
	e.NestedError = err
}

type ServiceErrorNotImplemented struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorNotImplemented) Error() string {
	return e.ErrorText
}

func (e *ServiceErrorNotImplemented) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorNotImplemented) SetNestedError(err error) {
	e.NestedError = err
}

type ServiceErrorInternalServerError struct {
	ErrorText   string
	NestedError error
}

func (e *ServiceErrorInternalServerError) SetErrorText(text string) {
	e.ErrorText = text
}

func (e *ServiceErrorInternalServerError) SetNestedError(err error) {
	e.NestedError = err
}

func (e *ServiceErrorInternalServerError) Error() string {
	return e.ErrorText
}

func newServiceError(nestedBackendError error, customMessage string, defaultMessage string, specificServiceError ServiceError) error {
	if nestedBackendError == nil && customMessage == "" {
		specificServiceError.SetErrorText(defaultMessage)
		specificServiceError.SetNestedError(nil)
		return eris.Wrap(specificServiceError, defaultMessage)
	} else if nestedBackendError != nil && customMessage == "" {
		specificServiceError.SetErrorText(nestedBackendError.Error())
		specificServiceError.SetNestedError(nestedBackendError)
		return eris.Wrap(specificServiceError, defaultMessage)
	} else if nestedBackendError == nil && customMessage != "" {
		specificServiceError.SetErrorText(defaultMessage)
		specificServiceError.SetNestedError(nil)
		return eris.Wrap(specificServiceError, customMessage)
	} else /* nestedBackendError != nil && customMessage != "" */ {
		specificServiceError.SetErrorText(nestedBackendError.Error())
		specificServiceError.SetNestedError(nestedBackendError)
		return eris.Wrap(specificServiceError, customMessage)
	}
}

func NewServiceErrorInternalServerError(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_INTERNAL_SERVER_ERROR", &ServiceErrorInternalServerError{})
}

func NewServiceErrorNotFound(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_NOT_FOUND", &ServiceErrorNotFound{})
}
func NewServiceErrorBadRequest(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_BAD_REQUEST", &ServiceErrorBadRequest{})
}

func NewServiceErrorNotImplemented(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_NOT_IMPLEMENTED", &ServiceErrorNotImplemented{})
}

func NewServiceErrorUnauthorized(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_UNAUTHORIZED", &ServiceErrorUnauthorized{})
}

func NewServiceErrorForbidden(nestedBackendError error, customMessage string) error {
	return newServiceError(nestedBackendError, customMessage, "SERVICE_ERROR_FORBIDDEN", &ServiceErrorForbidden{})
}
