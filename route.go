package vial

import (
    "fmt"
    "regexp"

    "github.com/pkg/errors"
)

// parameterRegex matches the surrounding declaration for a URL parameter.
var parameterRegex = regexp.MustCompile(`<([^>]+)>`)
var variableParse = regexp.MustCompile(
    `<(?:(?P<type>[^:]+):){0,1}(?P<name>[^>]+)>`,
)

// baseRegex grabs the URL up until the first path parameter.
var baseRegex = regexp.MustCompile(`([^<]+)`)

// uuidRegex matches a UUID.
var uuidRegex = `(?P<%s>[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?` +
    `[0-9a-fA-F]{4}-?[0-9a-fA-F]{12})`

// Route is a broken-down version of a route string.
type Route struct {
    original string
    matcher  *regexp.Regexp
    Base     string
}

// Matches will check if a url matches the route definition.
func (r Route) Matches(url string) bool {
    return r.matcher.MatchString(url)
}

// PathParams will take a url and parse the variables according to
// the route definition.
func (r Route) PathParams(url string) map[string]string {
    pathParamNameValueMap := getMappedValues(r.matcher, url)
    return pathParamNameValueMap
}

// getMappedValues matches the url and extracts the defined parameters.
func getMappedValues(regex *regexp.Regexp, input string) map[string]string {
    subMatch := regex.FindStringSubmatch(input)
    matchNames := regex.SubexpNames()
    subMatch, matchNames = subMatch[1:], matchNames[1:]
    mappedValues := make(map[string]string, len(subMatch))
    for i := range matchNames {
        mappedValues[matchNames[i]] = subMatch[i]
    }
    return mappedValues
}

func parseURLMatch(match string) string {
    variableMatches := variableParse.FindStringSubmatch(match)
    var newString string
    if len(variableMatches) == 3 {
        variableType := variableMatches[1]
        switch variableType {
        case "string":
            newString = `(?P<%s>[A-Za-z]+)`
        case "integer":
            newString = `(?P<%s>[0-9]+)`
        case "uuid":
            newString = uuidRegex
        case "":
            newString = `(?P<%s>[^\/\\]+)`
        default:
            panic(errors.Errorf("Unknown variable type: %s", variableType))
        }
        newString = fmt.Sprintf(newString, variableMatches[2])
    } else {
        panic(errors.Errorf("Incorrect format: %s", match))

    }
    return newString
}

// ParseRoute parses a route string with path param variable matchers into a
// Route struct.
func ParseRoute(route string) (newRoute Route, err error) {
    defer func() {
        // Because we can't error in ReplaceAllStringFunc it panics, so we
        // catch it here and bubble it.
        if r := recover(); r != nil {
            if recErr, ok := r.(error); ok {
                err = errors.Wrap(recErr, "Panic while parsing route string")
            }
        }
    }()

    matcher := parameterRegex.ReplaceAllStringFunc(route, parseURLMatch) + "$"
    matcherRegexp, newErr := regexp.Compile(matcher)
    if newErr != nil {
        err = errors.Wrapf(err, "Bad URL regex '%s'", matcher)
        return
    }
    baseURL := baseRegex.FindString(route)

    newRoute = Route{
        original: route,
        matcher:  matcherRegexp,
        Base:     baseURL,
    }

    return
}
