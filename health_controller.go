package vial

import (
    "daihasso.net/library/vial/responses"
)

// defaultHealthController returns a dummy value of true.
type defaultHealthController struct {}

// Get will return a simple true for now.
func (*defaultHealthController) Get(
    transactor *Transactor,
) responses.Data {
    healthy := true
    transactor.Logger.Debug("Health check.").
        With("requestor", transactor.Request.RemoteAddr).
        With("healthy", healthy).
        Send()

    return transactor.Respond(
        200,
        responses.Body(
            map[string]bool{
                "healthy": healthy,
            },
        ),
    )
}
