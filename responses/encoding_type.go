package responses

import (
)

// EncodingType is a way of encoding a response data.
type EncodingType int

const (
    UnsetEncoding EncodingType = iota
    JSONEncoding
    TextPlainEncoding
)
