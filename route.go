package vial

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/pkg/errors"
)

// parameterRegex matches the surrounding declaration for a URL parameter.
var parameterRegex = regexp.MustCompile(`<([^>]+)>`)
var variableParse = regexp.MustCompile(
    `<(?:(?P<type>[^:]+):){0,1}(?P<name>[^>]+)>`,
)

// baseRegex grabs the URL up until the first path parameter.
var baseRegex = regexp.MustCompile(`([^<]+)`)

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
func (r Route) PathParams(url string) (PathParams, error) {
    pathParamNameValueMap, err := getMappedValues(r.matcher, url)
    if err != nil {
        return nil, errors.Wrap(
            err, "Error while getting path parameter values",
        )
    }
    return pathParamNameValueMap, nil
}

func coerceType(key, typ, stringVal string) (interface{}, error) {
    pathParamMatcher, ok := GetPathParamMatcher(typ)
    if !ok {
        return nil, errors.Errorf(
            "No PathParamMatcher found for type '%s'", typ,
        )
    }
    val, err := pathParamMatcher.Coercer(stringVal)
    if err != nil {
        return nil, errors.Wrap(
            err, "Error while converting path parameter '%s' to type '%s'",
        )
    }

    return val, nil
}

// getMappedValues matches the url and extracts the defined parameters.
func getMappedValues(regex *regexp.Regexp, input string) (PathParams, error) {
    subMatch := regex.FindStringSubmatch(input)
    matchNames := regex.SubexpNames()
    subMatch, matchNames = subMatch[1:], matchNames[1:]
    mappedValues := make(PathParams, len(subMatch))
    for i := range matchNames {
        keyParts := strings.SplitN(matchNames[i], "_", 2)
        typ, key := keyParts[0], keyParts[1]
        val, err := coerceType(key, typ, subMatch[i])
        if err != nil {
            return nil, err
        }
        mappedValues[key] = val
    }
    return mappedValues, nil
}

func parseURLMatch(match string) string {
    variableMatches := variableParse.FindStringSubmatch(match)
    finalRegexTemplate := `(?P<%s_%s>%s)`
    var finalRegex string
    if len(variableMatches) == 3 {
        variableType := strings.ToLower(variableMatches[1])
        pathParamMatcher, ok := GetPathParamMatcher(variableType)
        if !ok {
            panic(errors.Errorf("Unknown variable type: %s", variableType))
        }

        finalRegex = fmt.Sprintf(
            finalRegexTemplate,
            pathParamMatcher.prefix(),
            variableMatches[2],
            pathParamMatcher.RegexString,
        )
    } else {
        panic(errors.Errorf("Incorrect format: %s", match))

    }

    return finalRegex
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

    if route[0] != '/' {
        route = "/" + route
    }

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
