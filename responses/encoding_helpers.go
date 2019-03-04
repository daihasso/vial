package responses

import (
)

var defaultEncodingContentTypeMap = map[EncodingType]ContentType{
    JSONEncoding: JSONContentType,
    TextPlainEncoding: TextPlainUTF8ContentType,
}

// DefaultContentTypeForEncoding will return the actual encoding-type string
// for use in a Content-Type header.
func DefaultContentTypeForEncoding(encodingType EncodingType) string {
    return string(defaultEncodingContentTypeMap[encodingType])
}
