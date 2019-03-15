package vial

import (
    "context"
)

// TransctorOption is an option for a Transactor.
type TransactorOption func(*Transactor) error

// WithExistingContext creates a transactor using the existing context which
// overrides the default `request.Context()`.
func WithExistingContext(existingCtx context.Context) TransactorOption {
    return func(transactor *Transactor) error {
        transactor.ChangeContext(existingCtx)

        return nil
    }
}
