package responses

import (
)

// AdditionalAttribute is an attribute added to a builder either after
// creation or on creation.
type AdditionalAttribute func(*Builder) error

// Body sets the body to the value(s) provided. If multiple values are provided
// a top-level array is assumed.
func Body(data ...interface{}) AdditionalAttribute {
    return func(rb *Builder) error {
        rb.SetBody(data...)
        return nil
    }
}

// Headers will merge provided headers with the existing headers. Keys that
// exist both in the existing headers and provided headers will use the value
// provided by the new headers over the existing headers.
func Headers(headers map[string][]string) AdditionalAttribute {
    return func(rb *Builder) error {
        rb.SetHeaders(headers)
        return nil
    }
}

// AddHeader adds a new value to the headers. If the header key exists it
// will add the provided header value(s) after the existing values.
func AddHeader(key, value string, otherValues ...string) AdditionalAttribute {
    return func(rb *Builder) error {
        rb.AddHeader(key, value, otherValues...)
        return nil
    }
}

// Status will set the status code for the response.
func Status(statusCode int) AdditionalAttribute {
    return func(rb *Builder) error {
        rb.SetStatus(statusCode)
        return nil
    }
}

func Encoding(encodingType EncodingType) AdditionalAttribute {
    return func(rb *Builder) error {
        rb.SetEncoding(encodingType)
        return nil
    }
}
