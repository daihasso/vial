package vial

import(
    "context"
    "fmt"
    "reflect"
    "strings"

    "daihasso.net/library/vial/responses"
)

// A RouteControllerMethod is a function which takes a context and a transactor
// and handles the request.
type RouteControllerMethod func(context.Context, *Transactor) responses.Data

// A RouteControllerMethodMinimal is a function which takes only a transactor
// and handles the request.
type RouteControllerMethodMinimal func(*Transactor) responses.Data

// A RouteController is a controller for a given route that defines the methods
// it responds to.
type RouteController interface{}

// A RouteControllerFunc is a functional controller for a given route & method
// that defines the methods it responds to.
type RouteControllerFunc interface{}

var routeControllerMethodType = reflect.TypeOf(
    func(context.Context, *Transactor) responses.Data{
        return responses.Data{}
    },
)
var routeControllerMethodMinimalType = reflect.TypeOf(
    func(*Transactor) responses.Data{
        return responses.Data{}
    },
)
var validRouteControllerTypes = map[reflect.Type]bool{
    routeControllerMethodType: true,
    routeControllerMethodMinimalType: true,
}
var validControllerFields = []string{
    "Post", "Get", "Put", "Patch", "Delete", "Head", "Options",
}

func validMethodsMessage() string {
    // TODO: This could probably be cached and only done once
    //       if that matters.
    var validMethods []string
    for methodType, _ := range validRouteControllerTypes {
        validMethods = append(
            validMethods, methodType.String(),
        )
    }
    return strings.Join(validMethods, "\nOR\n")
}

// RouteControllerIsValid checks a RouteController to see if it meets some of
// of the expectations of a RouteController.
func RouteControllerIsValid(rc RouteController) []string {
    var allErrMessages []string
    rcVal := reflect.ValueOf(rc)
    rcValRoot := rcVal
    for rcValRoot.Kind() == reflect.Ptr {
        rcValRoot = rcValRoot.Elem()
    }

    if rcValRoot.Kind() == reflect.Struct {
        // This might require a little explanation:
        // If there's a field with a valid controller name that matches one of
        // the approved types then this RouteController is valid, otherwise
        // it's not really a RouteController.
        // NOTE: Technically since we're using FieldByName and the code for
        //       that as of 1.28.19 iterates through all the fields to find a
        //       matching field this would cost something like:
        //         `O(x * y)`
        //       Where x is the number of fields on RouteController and y is
        //       the number of valid fields (defined by validControllerFields)
        //       Alternatively we could iterate all the fields manually and
        //       check their names while we iterate this would yield us:
        //         `O(x)`
        //       but that solution is a little less clean and I don't think the
        //       performance gain is worth it at this time.
        for _, methodName := range validControllerFields {
            if method := rcVal.MethodByName(methodName); method.IsValid() {
                if _, ok := validRouteControllerTypes[method.Type()]; !ok {
                    allErrMessages = append(
                        allErrMessages,
                        fmt.Sprintf(
                            "Method '%s' had signature:\n%s\nbut expected " +
                                "one of:\n%s",
                            methodName,
                            method.Type().String(),
                            validMethodsMessage(),
                        ),
                    )
                }
            }
        }
    } else if rcValRoot.Kind() == reflect.Func {
        if mcf, ok := rc.(methodControllerFunc); ok {
            if _, _, err := mcf(); err != nil {
                allErrMessages = append(
                    allErrMessages,
                    fmt.Sprintf(
                        "RouteControllerFunc provided to MethodFunc" +
                            " invalid:\n%s",
                        err.Error(),
                    ),
                )
            }
        } else {
            allErrMessages = append(
                allErrMessages,
                "Functional controllers should be wrapped in a MethodFunc" +
                    " call to tie it to a specific request method",
            )
        }
    }

    return allErrMessages
}
