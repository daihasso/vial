package vial

import (
    "fmt"
    "regexp"

    "github.com/pkg/errors"
    "github.com/google/uuid"
)

var (
    wrongPathParamType = "Path parameter is not of type %s"
    pathParamDoesNotExist = "Path parameter with key '%s' does not exist"
)

var (
    wrongPathParamTypeRegex = regexp.MustCompile(
        fmt.Sprintf(wrongPathParamType, `\w+`),
    )
    pathParamDoesNotExistRegex = regexp.MustCompile(
        fmt.Sprintf(pathParamDoesNotExist, `[^']*`),
    )
)

// IsWrongPathParamType checks if an error is caused by the type requested not
// being the correct type.
func IsWrongPathParamType(err error) bool {
    return wrongPathParamTypeRegex.MatchString(err.Error())
}

// IsPathParamDoesNotExist checks if an error is caused by the param not
// existing.
func IsPathParamDoesNotExist(err error) bool {
    return pathParamDoesNotExistRegex.MatchString(err.Error())
}

// PathParams is a convenience wrapper around a map holding coerced path
// values.
type PathParams map[string]interface{}

// String retrieves and coerces a string.
func (self PathParams) String(key string) (string, error) {
    if in, ok := self[key]; ok {
        if s, ok := in.(string); ok {
            return s, nil
        }

        return "", errors.New(fmt.Sprintf(wrongPathParamType, "string"))
    }

    return "", errors.New(fmt.Sprintf(pathParamDoesNotExist, key))
}

// Float retrieves and coerces a float.
func (self PathParams) Float(key string) (float64, error) {
    if in, ok := self[key]; ok {
        if f, ok := in.(float64); ok {
            return f, nil
        }

        return 0, errors.New(fmt.Sprintf(wrongPathParamType, "float"))
    }

    return 0, errors.New(fmt.Sprintf(pathParamDoesNotExist, key))
}

// Int retrieves and coerces a int.
func (self PathParams) Int(key string) (int, error) {
    if in, ok := self[key]; ok {
        if i, ok := in.(int); ok {
            return i, nil
        }

        return 0, errors.New(fmt.Sprintf(wrongPathParamType, "int"))
    }

    return 0, errors.New(fmt.Sprintf(pathParamDoesNotExist, key))
}

// UUID retrieves and coerces a UUID.
func (self PathParams) UUID(key string) (uuid.UUID, error) {
    if in, ok := self[key]; ok {
        if u, ok := in.(uuid.UUID); ok {
            return u, nil
        }

        return uuid.UUID{}, errors.New(fmt.Sprintf(wrongPathParamType, "UUID"))
    }

    return uuid.UUID{}, errors.New(fmt.Sprintf(pathParamDoesNotExist, key))
}
