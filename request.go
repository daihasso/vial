package vial

import (
    "context"
    "net/http"

    "github.com/google/uuid"
)

// InboundRequest is a (thin) wrapper around http.Request.
type InboundRequest struct {
    http.Request
    PathParams PathParams
}


// WithContext returns a new InboundRequest with the provided context.
func (r InboundRequest) WithContext(ctx context.Context) *InboundRequest {
    return &InboundRequest{
        Request: *r.Request.WithContext(ctx),
        PathParams: r.PathParams,
    }
}

// PathString returns a variable from the path matching the provided key as
// a string.
func (r InboundRequest) PathString(key string) (string, bool) {
    s, err := r.PathParams.String(key)
    return s, err == nil
}

// PathInt returns a path variable as a int or an error if it could not be
// found or could not be converted.
func (r *InboundRequest) PathFloat(key string) (float64, bool) {
    f, err := r.PathParams.Float(key)
    return f, err == nil
}

// PathInt returns a path variable as a int or an error if it could not be
// found or could not be converted.
func (r *InboundRequest) PathInt(key string) (int, bool) {
    i, err := r.PathParams.Int(key)
    return i, err == nil
}

// PathUUID returns a path variable as a UUID or an error if it could not be
// found or could not be converted.
func (r *InboundRequest) PathUUID(key string) (uuid.UUID, bool) {
    uuid, err := r.PathParams.UUID(key)

    return uuid, err == nil
}

// QueryParams returns the raw query params map.
func (r InboundRequest) QueryParams() map[string][]string {
    return r.URL.Query()
}

// QueryParamMultiple returns a parameter from the query parameters as a slice
// of all the results.
func (r *InboundRequest) QueryParamMultiple(key string) ([]string, bool) {
    value, ok := r.URL.Query()[key]

    return value, ok
}

// QueryParam returns a parameter from the query parameters using the
// first result if there are many.
func (r *InboundRequest) QueryParam(key string) (string, bool) {
    values, ok := r.URL.Query()[key]

    if !ok || len(values) == 0 {
        return "", false
    }

    return values[0], true
}

// NewInboundRequest gets a request based on an existing HTTP request with
// path parameters included.
func NewInboundRequest(
    baseRequest *http.Request,
    pathParams PathParams,
) *InboundRequest {
    req := &InboundRequest{
        Request: *baseRequest,
        PathParams: pathParams,
    }
    return req
}
