package vial

import (
    "net/http"
    "strconv"

    "github.com/google/uuid"
)

// RequestDirection is the direction that the request is going inbound
// (received from a client) or outbound (to be sent from a client).
type RequestDirection int

// Enumerate the RequestDirections
const (
    _ RequestDirection = iota
    InboundRequest
    OutboundRequest
)


// Request is a (thin) wrapper around http.Request.
type Request struct {
    http.Request
    PathParams map[string]string
    direction RequestDirection
}

// PathString returns a variable from the path matching the provided key as
// a string.
func (r Request) PathString(key string) (string, bool) {
    val, ok := r.PathParams[key]

    return val, ok
}

// PathInt returns a path variable as a int or an error if it could not be
// found or could not be converted.
func (r *Request) PathInt(key string) (int, bool) {
    if value, ok := r.PathParams[key]; ok {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, false
        }

        return intValue, true
    }

    return 0, false
}

// PathUUID returns a path variable as a UUID or an error if it could not be
// found or could not be converted.
func (r *Request) PathUUID(key string) (uuid.UUID, bool) {
    var uuidValue uuid.UUID
    var err error
    value, ok := r.PathParams[key]
    if !ok {
        return uuidValue, false
    }

    uuidValue, err = uuid.Parse(value)
    if err != nil {
        return uuidValue, false
    }

    return uuidValue, true
}

// QueryParams returns the raw query params map.
func (r Request) QueryParams() map[string][]string {
    return r.URL.Query()
}

// QueryParamMultiple returns a parameter from the query parameters as a slice
// of all the results.
func (r *Request) QueryParamMultiple(key string) ([]string, bool) {
    value, ok := r.URL.Query()[key]

    return value, ok
}

// QueryParam returns a parameter from the query parameters using the
// first result if there are many.
func (r *Request) QueryParam(key string) (string, bool) {
    values, ok := r.URL.Query()[key]

    if !ok || len(values) == 0 {
        return "", false
    }

    return values[0], true
}

// NewServerRequest gets a request based on an existing HTTP request with
// path parameters included.
func NewServerRequest(
    baseRequest *http.Request,
    pathParams map[string]string,
) *Request {
    req := &Request{
        Request: *baseRequest,
        PathParams: pathParams,
        direction: InboundRequest,
    }
    return req
}
