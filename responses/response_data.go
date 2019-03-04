package responses

import (
    "net/http"

    "github.com/pkg/errors"
)

// Data is a representation of the route transaction's response.
type Data struct {
    Headers map[string][]string
    Body []byte
    StatusCode int

    // unexpectedError is an error that's incidental (ex: an error that
    // happened while forming the ResponseData itself). It will trigger an
    // Internal Service Error (500).
    unexpectedError error
}

// Write writes the data inside of a Data struct into a response writer.
func (self Data) Write(w http.ResponseWriter) error {
    for key, value := range self.Headers {
        w.Header()[key] = value
    }
    w.WriteHeader(self.StatusCode)
    _, err := w.Write(self.Body)

    return errors.Wrap(err, "Error while writing response data to response")
}

// Error returns and unexpected errors that may have happened.
func (self Data) Error() error {
    return self.unexpectedError
}

// ErrorResponse returns a response Data from an error. This is for fatal
// abortions, generally avoid this in favor of returning a neterr.CodedError
// with abort instead for better client experience.
func ErrorResponse(err error) Data {
    return Data{unexpectedError: err}
}
