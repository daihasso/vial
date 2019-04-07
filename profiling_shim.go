// +build !pprof

package vial

import (
    "net/http"
)

func handleProfiling(muxer *http.ServeMux) {
    // NOOP because we're not enabling profiling.
}
