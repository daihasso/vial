package vial

import (
    "context"

    "github.com/pkg/errors"
    "github.com/google/uuid"
)

// This is the context base-key for all vial context keys.
var ContextKeyBase = "vial."

// This is the key that the Sequence ID is stored under in the context.
var SequenceIdContextKey = ContextKey("sequence_id")

// This is the key that the server is stored under.
var ServerContextKey = ContextKey("server")

// This is the key that the server's logger is stored under.
var ServerLoggerContextKey = ContextKey("server.logger")

// This is the key that the transactor will live under in the request context.
var TransactorContextKey = ContextKey("transactor")

// This is the key that the Sequence ID is stored under in the context.
var RequestIdContextKey = ContextKey("request_id")

// ContextKey is a helper for generating a context key prefixed for vial.
func ContextKey(key string) string {
    return ContextKeyBase + key
}

// ContextWithSequenceId creates a new context with the provided Sequence ID
// inserted into it.
func ContextWithSequenceId(
    ctx context.Context, sequenceId uuid.UUID,
) context.Context {
    return context.WithValue(ctx, SequenceIdContextKey, sequenceId.String())
}

// ContextSequenceId tries to retrieve a Sequence ID stored in the provided
// context and returns it if it exists.
func ContextSequenceId(ctx context.Context) (*uuid.UUID, error) {
    sequenceIdIn := ctx.Value(SequenceIdContextKey)
    if sequenceIdIn == nil {
        return nil, errors.New(
            "SequenceId not in context. Is this call made via vial?",
        )
    }

    sequenceIdString, ok := sequenceIdIn.(string)
    if !ok {
        return nil, errors.New(
            "Value of SequenceId in context is not a string",
        )
    }

    sequenceId, err := uuid.Parse(sequenceIdString)
    if err != nil {
        return nil, errors.Wrap(
            err,
            "Value of SequenceId in context is not a proper UUID",
        )
    }

    return &sequenceId, nil
}

func ContextRequestId(ctx context.Context) (*uuid.UUID, error) {
    requestIdIn := ctx.Value(RequestIdContextKey)
    if requestIdIn == nil {
        return nil, errors.New(
            "RequestId not in context. Is this call made via vial?",
        )
    }

    requestIdString, ok := requestIdIn.(string)
    if !ok {
        return nil, errors.New(
            "Value of RequestId in context is not a string",
        )
    }

    requestId, err := uuid.Parse(requestIdString)
    if err != nil {
        return nil, errors.Wrap(
            err,
            "Value of RequestId in context is not a proper UUID",
        )
    }

    return &requestId, nil
}
