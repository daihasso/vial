// Package neterr seeks to unify error formatting for API consumers. All
// CodedErrors have a code (int) that is unique to that error within your API
// and a message (string) that describes what the error means.
package neterr

import (
    "encoding/json"
)

// CodedError is an error with a code and a message. It is used to standardize
// error responses for a uniform API consumer experience.
// Some special errors are Vial-specific errors which are noted by having a
// true return value from the IsVialError method call.
type CodedError struct {
    code int
    message string
    isVialError bool
}

type marshalableCodedError struct{
    Code int `json:"code"`
    Message string `json:"message"`
    IsVialError bool `json:"vial_error,omitempty"`
}

func (self *CodedError) UnmarshalJSON(data []byte) error {
    tempStruct := marshalableCodedError{
        Code: self.code,
        Message: self.message,
        IsVialError: self.isVialError,
    }
    err := json.Unmarshal(data, &tempStruct)
    if err != nil {
        return err
    }

    self.code = tempStruct.Code
    self.message = tempStruct.Message
    self.isVialError = tempStruct.IsVialError

    return nil
}

func (self CodedError) MarshalJSON() ([]byte, error) {
    return json.Marshal(
        marshalableCodedError{
            Code: self.code,
            Message: self.message,
            IsVialError: self.isVialError,
        },
    )
}

// Message returns the message attached to the CodedError.
func (self *CodedError) Message() string {
    return self.message
}

// Code returns the error code for the CodedError.
func (self *CodedError) Code() int {
    return self.code
}

// IsVialError specifies whether this error is a framework-level error or
// not.
func (self *CodedError) IsVialError() bool {
    return self.isVialError
}


// NewCodedError creates a new CodedError.
func NewCodedError(code int, message string) CodedError {
    return CodedError{
        code: code,
        message: message,
        isVialError: false,
    }
}

// CodedErrorFromError is a helper that takes a go error and formats it into a
// proper CodedError.
func CodedErrorFromError(code int, err error) CodedError {
    return CodedError{
        code: code,
        message: err.Error(),
        isVialError: false,
    }
}
