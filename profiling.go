// +build pprof

package vial

import (
    "net/http"
    pprof "net/http/pprof"
)

func handleProfiling(muxer *http.ServeMux) {
    muxer.HandleFunc("/debug/pprof/", pprof.Index)
    muxer.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    muxer.HandleFunc("/debug/pprof/profile", pprof.Profile)
    muxer.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    muxer.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
