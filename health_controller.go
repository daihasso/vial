package vial

import (
    "github.com/daihasso/slogging"

    "daihasso.net/library/vial/responses"
)

// defaultHealthController returns a dummy value of true.
type defaultHealthController struct {}

// Get will return a simple true for now.
func (*defaultHealthController) Get(
    transactor *Transactor,
) responses.Data {
    healthy := true
    transactor.Logger.Debug("Health check.", logging.Extras{
        "requestor": transactor.Request.RemoteAddr,
        "healthy": healthy,
    })

    return transactor.Respond(
        200,
        responses.Body(
            map[string]bool{
                "healthy": healthy,
            },
        ),
    )
}
