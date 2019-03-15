package vial


import (
    "context"

    "github.com/pkg/errors"

    "github.com/daihasso/vial/responses"
)

// PreMiddleWare is a function that's run against a request before the route's
// primary function is performed.
// Returning a Data struct will short-circut the response and use this response
// instead of processing the normal handler.
// Returning a context will make that the new context for future middleware or
// handlers.
type PreMiddleWare func(context.Context, *Transactor) (
    *responses.Data, *context.Context, error,
)

// PostMiddleWare is a function that's run against a request after the route's
// primary function is performed.
type PostMiddleWare func(
    context.Context, *Transactor, responses.Data,
) (*responses.Data, error)

func DefaultEncryptionHeadersMiddleware() PreMiddleWare {
    return func(
        _ context.Context, transactor *Transactor,
    ) (*responses.Data, *context.Context, error) {
        err := transactor.AddHeader(
            "Strict-Transport-Security",
            "max-age=63072000; includeSubDomains",
        )

        if err != nil {
            return nil, nil, errors.Wrap(
                err,
                "Error while adding default encryption header to response",
            )
        }

        return nil, nil, nil
    }
}
