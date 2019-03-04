package vial

import (
    "sort"
    "strings"
    "net/http"

    "daihasso.net/library/vial/responses"
    "daihasso.net/library/vial/neterr"
)

// DefaultOptions iterates a routes available methods and returns them in the
// header.
func DefaultOptions(
    server Server,
    transactor *Transactor,
    methods []*RouteControllerHelper,
) responses.Data {
    var methodStrings []string
    seenMethods := map[RequestMethod]bool{
        MethodOPTIONS: true,
    }
    for _, rch := range(methods) {
        for _, method := range rch.AllMethods() {
            if _, ok := seenMethods[method]; !ok {
                methodStrings = append(methodStrings, method.String())
                seenMethods[method] = true
            }
        }
    }
    sort.Strings(methodStrings)
    methodStrings = append([]string{MethodOPTIONS.String()}, methodStrings...)
    err := transactor.SetHeader(
        "Access-Control-Allow-Methods",
        strings.Join(methodStrings, ", "),
    )
    if err != nil {
        return transactor.Abort(
            http.StatusInternalServerError,
            neterr.DefaultOptionsHeaderSetError,
            neterr.CodedErrorFromError(0, err),
        )
    }

    return transactor.Respond(http.StatusOK)
}
