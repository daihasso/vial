// Package responses handles common response logic and helpers for constructing
// and managing responses.
package responses

import (
    "encoding/json"
    "fmt"
    "context"
    "net/http"
    "net/textproto"

    "github.com/pkg/errors"

    "github.com/daihasso/vial/neterr"
)

// Builder is a helper struct for constructing a Data struct for a transaction
// response.
type Builder struct {
    ctx context.Context
    isArray bool
    bodyData interface{}
    headers map[string][]string
    statusCode int
    encodingType EncodingType
    contentType string
}

func (self *Builder) finalizeJSONBody() ([]byte, error) {
    bytes, err := json.Marshal(self.bodyData)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

func (self *Builder) applyAdditionals(
    additionals []AdditionalAttribute,
) error {
    var err error
    for _, additional := range additionals {
        err = additional(self)
        if err != nil {
            return errors.Wrap(
                err, "Error while running addditionals on Builder",
            )
        }
    }

    return nil
}

func (self Builder) prepare(
    additionals []AdditionalAttribute,
) *Data {
    err := self.applyAdditionals(additionals)
    if err != nil {
        return &Data{
            unexpectedError: errors.Wrap(
                err, "Error while adding finishing addditionals to Builder",
            ),
        }
    }

    httpHeaders := http.Header{}
    for key, value := range self.headers {
        properKey := textproto.CanonicalMIMEHeaderKey(key)
        httpHeaders[properKey] = value
    }

    if self.contentType != "" {
        httpHeaders.Set(ContentTypeHeader, self.contentType)
    } else if _, ok := httpHeaders[ContentTypeHeader]; !ok {
        httpHeaders.Set(
            ContentTypeHeader,
            DefaultContentTypeForEncoding(self.encodingType),
        )
    }

    var body []byte
    // NOTE: If we get a straight string or bytes just treat it as-is instead
    //       of marshalling it. This might be controversial but I think the
    //       number of people providing pre-processed/marshalled bodies to the
    //       builder will vastly outweigh the legitimate needs for a top-level
    //       marshalled string/bytes.
    if bodyString, ok := self.bodyData.(string); ok {
        body = []byte(bodyString)
    } else if bodyBytes, ok := self.bodyData.([]byte); ok {
        body = bodyBytes
    } else {
        switch self.encodingType {
        case JSONEncoding:
            body, err = self.finalizeJSONBody()
        case TextPlainEncoding:
            body = []byte(fmt.Sprint(self.bodyData))
        case UnsetEncoding:
            err = errors.New(
                "Encoding type not set, unsure how to marshal content.",
            )
        default:
            err = errors.Errorf(
                "Unknown encodingType: %#+v", self.encodingType,
            )
        }
    }

    return &Data{
        Headers: httpHeaders,
        Body: body,
        StatusCode: self.statusCode,
        unexpectedError: err,
    }
}

// SetBody sets the body to the provided data. If multiple values are provided,
// a top-level array is assumed.
func (self *Builder) SetBody(data ...interface{}) {
    if len(data) > 1 {
        self.isArray = true
        self.bodyData = data
    } else if len(data) > 0 {
        self.isArray = false
        self.bodyData = data[0]
    }
}

// ReplaceHeaders replaces all the existing headers with the provided headers.
func (self *Builder) ReplaceHeaders(headers map[string][]string) {
    self.headers = make(map[string][]string, len(headers))

    for key, values := range headers {
        properKey := textproto.CanonicalMIMEHeaderKey(key)
        self.headers[properKey] = values
    }
}

// SetHeaders will merge provided headers with the existing headers. Keys that
// exist both in the existing headers and provided headers will use the value
// provided by the new headers over the existing headers.
func (self *Builder) SetHeaders(headers map[string][]string) {
    for key, values := range headers {
        properKey := textproto.CanonicalMIMEHeaderKey(key)
        self.headers[properKey] = values
    }
}

// AddHeaders will merge provided headers with the existing headers. Keys that
// exist both in the existing headers and provided headers will append the
// new values after the existing values.
func (self *Builder) AddHeaders(headers map[string][]string) {
    for key, values := range headers {
        properKey := textproto.CanonicalMIMEHeaderKey(key)
        if existingValues, ok := self.headers[properKey]; ok {
            self.headers[properKey] = append(existingValues, values...)
        } else {
            self.headers[properKey] = values
        }
    }
}

// SetHeader will set the header with provided key to the provided value(s). If
// a header with the key provided already exists it will be overriden by the
// value(s) provided.
func (self *Builder) SetHeader(key, value string, otherValues ...string) {
    self.headers[textproto.CanonicalMIMEHeaderKey(key)] = append(
        []string{value}, otherValues...,
    )
}

// AddHeader appends the value(s) provided to the header with the provided key.
// If no header with the provided key exists it is added with the provided
// values.
func (self *Builder) AddHeader(key, value string, otherValues ...string) {
    properKey := textproto.CanonicalMIMEHeaderKey(key)
    if existingValues, ok := self.headers[properKey]; ok {
        self.headers[properKey] = append(
            existingValues, append([]string{value}, otherValues...)...,
        )
    } else {
        self.headers[properKey] = append([]string{value}, otherValues...)
    }
}

// SetStatus sets the status for the response.
func (self *Builder) SetStatus(statusCode int) {
    self.statusCode = statusCode
}

// SetEncoding sets the encoding for the response.
func (self *Builder) SetEncoding(encoding EncodingType) {
    self.encodingType = encoding
}

// SetContentType sets the content type for the response.
func (self *Builder) SetContentType(contentType string) {
    self.contentType = contentType
}

// Abort generates a failure response.
func (self *Builder) Abort(
    statusCode int,
    codedError neterr.CodedError,
    otherErrors ...neterr.CodedError,
) Data {
    codedErrors := append([]neterr.CodedError{codedError}, otherErrors...)
    self.statusCode = statusCode
    // NOTE: This overrides the already set body, I think this is the best
    //       approach in this situation to prevent data leakage but it may be
    //       debatable.
    self.bodyData = map[string]interface{}{
        "errors": codedErrors,
    }

    return *self.prepare(nil)
}

// Finish finishes the builder generating a Data struct to return as a
// response.
func (self *Builder) Finish(additionals ...AdditionalAttribute) Data {
    return *self.prepare(additionals)
}

// NewBuilder generates a new builder to help generate a response Data struct.
func NewBuilder(
    ctx context.Context,
    encoding EncodingType,
    additionals ...AdditionalAttribute,
) (*Builder, error) {
    responseBuilder := new(Builder)
    responseBuilder.encodingType = encoding
    responseBuilder.ctx = ctx
    responseBuilder.headers = make(map[string][]string)

    err := responseBuilder.applyAdditionals(additionals)
    if err != nil {
        return nil, err
    }

    return responseBuilder, nil
}

// Abort generates an abort Data struct for responding with an error.
func Abort(
    ctx context.Context,
    encoding EncodingType,
    statusCode int,
    codedError neterr.CodedError,
    otherErrors ...neterr.CodedError,
) Data {
    codedErrors := append([]neterr.CodedError{codedError}, otherErrors...)
    builder, err := NewBuilder(
        ctx,
        encoding,
        Status(statusCode),
        Body(map[string]interface{}{
            "errors": codedErrors,
        }),
        Headers(map[string][]string{
            ContentTypeHeader: []string{JSONContentType},
        }),
    )
    if err != nil {
        return Data{unexpectedError: err}
    }

    return *builder.prepare(nil)
}

// Respond is a shortcut for generating a response Data struct via a builder
// all at once.
func Respond(
    ctx context.Context,
    encoding EncodingType,
    statusCode int,
    additionals ...AdditionalAttribute,
) Data {
    allAttributes := append(
        []AdditionalAttribute{Status(statusCode)}, additionals...,
    )
    builder, err := NewBuilder(ctx, encoding, allAttributes...)
    if err != nil {
        return Data{
            unexpectedError: errors.Wrap(
                err, "Error while creating new builder",
            ),
        }
    }
    return *builder.prepare(nil)
}
