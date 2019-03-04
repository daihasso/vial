package responses

import (
    "net/textproto"
)

// ContentTypeHeader is the Content-Type header ran through the STD library
// formatter.
var ContentTypeHeader = textproto.CanonicalMIMEHeaderKey("Content-Type")

// ContentType is a representation of a content type header value.
type ContentType string

const (
    TextPlainUTF8ContentType = "text/plain; charset=utf-8"
    JSONContentType = "application/json; charset=utf-8"
)
