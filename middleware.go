package vial


import (
    "context"

    "github.com/pkg/errors"

    "daihasso.net/library/vial/responses"
)

// PreMiddleWare is a function that's run against a request before the route's
// primary function is performed.
type PreMiddleWare func(context.Context, *Transactor) (*responses.Data, error)

// PostMiddleWare is a function that's run against a request after the route's
// primary function is performed.
type PostMiddleWare func(
    context.Context, *Transactor, responses.Data,
) (*responses.Data, error)

func DefaultEncryptionHeadersMiddleware() PreMiddleWare {
    return func(
        _ context.Context, transactor *Transactor,
    ) (*responses.Data, error) {
        err := transactor.AddHeader(
            "Strict-Transport-Security",
            "max-age=63072000; includeSubDomains",
        )

        if err != nil {
            return nil, errors.Wrap(
                err,
                "Error while adding default encryption header to response",
            )
        }

        return nil, nil
    }
}
