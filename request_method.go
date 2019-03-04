package vial

import (
	"strings"
	"net/http"
)

// RequestMethod is a method of request (i.e: GET).
type RequestMethod int

// Describes all the RequestMethods.
const (
	_ RequestMethod = iota
	MethodUnknown

	MethodGET
	MethodPUT
	MethodPOST
	MethodDELETE
	MethodOPTIONS
	MethodHEAD
	MethodPATCH
)

// RequestMethodFromString converts a string HTTP method to a RequestMethod
// enum.
func RequestMethodFromString(method string) RequestMethod {
	switch strings.ToUpper(method) {
	case http.MethodPost:
		return MethodPOST
	case http.MethodGet:
		return MethodGET
	case http.MethodPut:
		return MethodPUT
	case http.MethodPatch:
		return MethodPATCH
	case http.MethodDelete:
		return MethodDELETE
	case http.MethodHead:
		return MethodHEAD
	case http.MethodOptions:
		return MethodOPTIONS
	default:
		return MethodUnknown
	}
}

// RequestMethodFromString converts a string HTTP method to a RequestMethod
// enum.
func (self RequestMethod) String() string {
	switch self {
	case MethodPOST:
		return http.MethodPost
	case MethodGET:
		return http.MethodGet
	case MethodPUT:
		return http.MethodPut
	case MethodPATCH:
		return http.MethodPatch
	case MethodDELETE:
		return http.MethodDelete
	case MethodHEAD:
		return http.MethodHead
	case MethodOPTIONS:
		return http.MethodOptions
	default:
		return ""
	}
}
