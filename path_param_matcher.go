package vial

import (
    "sync"
    "strconv"
    "strings"

    "github.com/google/uuid"
)

var once sync.Once
var pathParamMatcherMutex = new(sync.RWMutex)
var pathParamMatchers = map[string]*PathParamMatcher{}


// PathParamCoercer is a function which takes a raw path param string value and
// converts it to the expected value.
type PathParamCoercer func(string) (interface{}, error)

// PathParamsMatcher is a matcher for a path parameter specified in the
// <type:var> part of a route.
type PathParamMatcher struct {
    // Identifiers is a list of all the aliases that this path parameter,
    // matches in a route definition such as <int:id> where int is an
    // identifier.
    Identifiers []string

    // RegexString is a string containing a regex that matches this type of
    // PathParam.
    RegexString string

    // Coercer as specified above takes a string value and returns it's
    // coerced value.
    Coercer PathParamCoercer
}

func (self PathParamMatcher) prefix() string {
    return strings.ToLower(self.Identifiers[0])
}

// AddPathParamMatcher adds a new path param matcher to the global registry.
func AddPathParamMatcher(newMatcher *PathParamMatcher) {
    pathParamMatcherMutex.Lock()
    defer pathParamMatcherMutex.Unlock()

    for _, identifier := range newMatcher.Identifiers {
        pathParamMatchers[strings.ToLower(identifier)] = newMatcher
    }
}

// GetPathParamMatcher retrieves a PathParamMatcher for a given identifier.
func GetPathParamMatcher(identifier string) (*PathParamMatcher, bool) {
    pathParamMatcherMutex.RLock()
    defer pathParamMatcherMutex.RUnlock()

    in, ok := pathParamMatchers[identifier]
    return in, ok
}

// StringPathParamMatcher is the most basic PathParamMatcher that simply
// matches any non-forward-slash character. It is also the default behaviour.
var StringPathParamMatcher = &PathParamMatcher{
    Identifiers: []string{"string", ""},
    RegexString: `[^\/\\]+`,
    Coercer: func(stringVal string) (interface{}, error) {
        return stringVal, nil
    },
}

// IntPathParamMatcher matches only whole integers.
var IntPathParamMatcher = &PathParamMatcher{
    Identifiers: []string{"int", "integer"},
    RegexString: `[0-9]+`,
    Coercer: func(stringVal string) (interface{}, error) {
        return strconv.Atoi(stringVal)
    },
}

// FloatPathParamMatcher matches only whole floats numbers. Floats are defined
// as having at least a decimal value.
var FloatPathParamMatcher = &PathParamMatcher{
    Identifiers: []string{"float"},
    RegexString: `[0-9]*\.[0-9]+`,
    Coercer: func(stringVal string) (interface{}, error) {
        return strconv.ParseFloat(stringVal, 64)
    },
}

// UUIDPathParamMatcher matches the generic UUID format and uses google's UUID
// library to convert the uuid string to a uuid.UUID.
var UUIDPathParamMatcher = &PathParamMatcher{
    Identifiers: []string{"uuid"},
    RegexString: `[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?` +
    `[0-9a-fA-F]{4}-?[0-9a-fA-F]{12}`,
    Coercer: func(stringVal string) (interface{}, error) {
        return uuid.Parse(stringVal)
    },
}

func init() {
    once.Do(func() {
        AddPathParamMatcher(StringPathParamMatcher)
        AddPathParamMatcher(IntPathParamMatcher)
        AddPathParamMatcher(FloatPathParamMatcher)
        AddPathParamMatcher(UUIDPathParamMatcher)
    })
}
