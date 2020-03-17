package vial

import (
    "context"

    "github.com/daihasso/slogging"

    "github.com/daihasso/vial/responses"
)

type RouteFunction func(transactor *Transactor) responses.Data

var validRouteActions = []string{
    "post", "get", "put", "patch", "delete", "head", "options",
}

func WithContext(
    controller func(context.Context, *Transactor) responses.Data,
) RouteFunction {
    return func(transactor *Transactor) responses.Data {
        return controller(transactor.Context(), transactor)
    }
}

func ContextOnly(
    controller func(context.Context) responses.Data,
) RouteFunction {
    return func(transactor *Transactor) responses.Data {
        return controller(transactor.Context())
    }
}

func (self *Server) addRoute(
    method, path string,
    callback RouteFunction,
    otherCallbacks ...RouteFunction,
) error {
    logging.Debug("Adding route callback(s).", logging.Extras{
        "http_method": method,
        "route_path": path,
        "callbacks": len(otherCallbacks) + 1,
    })

    // TODO: Stop falling back on the overly complicated AddController method.
    others := make([]RouteController, len(otherCallbacks))
    for i, cb := range otherCallbacks {
        others[i] = FuncHandler(method, (func(*Transactor) responses.Data)(cb))
    }
    return self.AddController(
        path,
        FuncHandler(method, (func(*Transactor) responses.Data)(callback)),
        others...,
    )
}

func (self *Server) Post(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("post", path, callback, otherCallbacks...)
}

func (self *Server) Get(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("get", path, callback, otherCallbacks...)
}

func (self *Server) Put(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("put", path, callback, otherCallbacks...)
}

func (self *Server) Patch(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("patch", path, callback, otherCallbacks...)
}

func (self *Server) Delete(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("delete", path, callback, otherCallbacks...)
}

func (self *Server) Head(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("head", path, callback, otherCallbacks...)
}

func (self *Server) Options(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    return self.addRoute("options", path, callback, otherCallbacks...)
}

func (self *Server) All(
    path string, callback RouteFunction, otherCallbacks ...RouteFunction,
) error {
    for _, method := range validRouteActions {
        err := self.addRoute(
            method, path, callback, otherCallbacks...,
        )
        if err != nil {
            return err
        }
    }

    return nil
}
