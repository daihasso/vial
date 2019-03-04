package vial

import (
)

type ServerChannelResponse struct {
    Type ServerChannelResponseType
    Error error
}


type ServerChannelResponseType int

const (
    UnknownChannelResponse ServerChannelResponseType = iota
    ServerStartChannelResponse
    ServerShutdownChannelResponse
    UnknownErrorChannelResponse
)
