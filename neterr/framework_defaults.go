package neterr

import (
)

func newVialError(code int, message string) CodedError {
    return CodedError{
        code: code,
        message: message,
        isVialError: true,
    }
}

// DefaultOptionsHeaderSetError occurs when there is a problem setting the
// headers in the DefaultOptions route.
var DefaultOptionsHeaderSetError = newVialError(
    4,
    "Error while settings headers in options route.",
)

// SwaggerNotFoundError is an error that occurs when the swagger does not
// exist.
var SwaggerNotFoundError = newVialError(
    3,
    "Couldn't find swagger file.",
)

// RouteNotSetupError is an error that occurs when a route was called that
// hasn't been setup.
var RouteNotSetupError = newVialError(
    2,
    "Route specified has not been setup.",
)

// MethodNotAllowedErrror is send when a method that is not setup is called on
// a route that exists.
var MethodNotAllowedError = newVialError(
    1,
    "Method not allowed.",
)
