package vial

import (
    "context"
    "fmt"
    "reflect"

    "github.com/pkg/errors"
    logging "github.com/daihasso/slogging"

    "daihasso.net/library/vial/responses"
)

// RouteControllerCaller provides a uniform function for calling any of the
// controller function variants.
type RouteControllerCaller func(context.Context, *Transactor) responses.Data

// methodControllerFunc is a function that returns the HTTP method with a func.
type methodControllerFunc func() (
    []RequestMethod, RouteControllerCaller, error,
)

func standardRouteControllerCaller(
    rcmm RouteControllerMethod,
) RouteControllerCaller {
    caller := func(
        ctx context.Context, transactor *Transactor,
    ) responses.Data {
        return rcmm(ctx, transactor)
    }

    return caller
}

func minimalRouteControllerCaller(
    rcmm RouteControllerMethodMinimal,
) RouteControllerCaller {
    caller := func(
        ctx context.Context, transactor *Transactor,
    ) responses.Data {
        return rcmm(transactor)
    }

    return caller
}

// RouteControllerHelper is a wrapper around a RouteController that stores its
// methods and the route it's attached to.
type RouteControllerHelper struct {
    route Route
    methodCallers map[RequestMethod]RouteControllerCaller
}

// AllMethods returns all the methods the RouteControllerCaller responds to.
func (self RouteControllerHelper) AllMethods() []RequestMethod {
    var methods []RequestMethod
    for method, _ := range self.methodCallers {
        methods = append(methods, method)
    }
    return methods
}

// RespondsToMethod checks if the RouteController this helper wraps responds to
// the given method.
func (self RouteControllerHelper) RespondsToMethod(m RequestMethod) bool {
    _, ok := self.methodCallers[m]
    return ok
}

// RespondsToMethodString checks if the RouteController this helper wraps
// responds to the given method string.
func (self RouteControllerHelper) RespondsToMethodString(
    methodString string,
) bool {
    m := RequestMethodFromString(methodString)
    _, ok := self.methodCallers[m]
    return ok
}

// ControllerFuncForMethod will return the correct function for a given HTTP
// method. This method expects you to only call methods that are defined. Use
// RespondsToMethod to check before calling this function.
func (self RouteControllerHelper) ControllerFuncForMethod(
    m RequestMethod,
) (RouteControllerCaller, bool) {
    caller, ok := self.methodCallers[m]
    return caller, ok
}

func wrapControllerMethod(in interface{}) (RouteControllerCaller, error) {
    switch v := in.(type) {
        case func(context.Context, *Transactor) responses.Data:
        return standardRouteControllerCaller(v), nil
        case func(*Transactor) responses.Data:
        return minimalRouteControllerCaller(v), nil
    }

    return nil, errors.Errorf(
        "Unknown controller function type '%T'", in,
    )
}

// FuncHandler wraps a functional method handler with the methods it responds
// to.
func FuncHandler(
    rcf RouteControllerFunc, method string, otherMethods ...string,
) methodControllerFunc {
    var reqMethods []RequestMethod
    methods := append([]string{method}, otherMethods...)
    for _, methodString := range methods {
        reqMethods = append(reqMethods, RequestMethodFromString(methodString))
    }

    rcWrap, err := wrapControllerMethod(rcf)

    return func() ([]RequestMethod, RouteControllerCaller, error) {
        if err != nil {
            return nil, nil, errors.Wrapf(
                err,
                "RouteControllerFunc provided had signature:\n%T but" +
                    " expected one of:\n%s\n",
                rcf,
                validMethodsMessage(),
            )
        }
        return reqMethods, rcWrap, nil
    }
}

// MethodsForRouteController gets all the methods that a RouteController has
// available.
// NOTE: This method implicitly makes the assumption that the first controllers
//       take presidence (i.e. If both controller 1 and controller 5 respond to
//       the GET method then controller 1 will be the controller chosen).
func MethodsForRouteController(
    path string,
    routeControllers ...RouteController,
) map[RequestMethod]RouteControllerCaller {
    methods := make(map[RequestMethod]RouteControllerCaller)

    for i, rc := range routeControllers {
        rcVal := reflect.ValueOf(rc)
        rcValRoot := rcVal
        for rcValRoot.Kind() == reflect.Ptr {
            rcValRoot = rcValRoot.Elem()
        }

        if rcValRoot.Kind() == reflect.Struct {
            for _, fieldName := range validControllerFields {
                if method := rcVal.MethodByName(fieldName); method.IsValid() {
                    in := method.Interface()
                    rcWrap, err := wrapControllerMethod(in)
                    if err != nil {
                        panic(err)
                    }
                    reqMethod := RequestMethodFromString(fieldName)
                    if _, ok := methods[reqMethod]; ok {
                        logging.Warn(fmt.Sprintf(
                            "RouteController #%d passed to " +
                                "AddController for route '%s' both " +
                                "respond to the %s method.",
                            i,
                            path,
                            reqMethod.String(),
                        )).Send()
                    }
                    methods[reqMethod] = rcWrap
                }
            }
        } else if mcf, ok := rc.(methodControllerFunc); ok {
            funcMethods, rcWrapper, _ := mcf()
            for _, method := range funcMethods {
                methods[method] = rcWrapper
            }
        }
    }

    return methods
}
