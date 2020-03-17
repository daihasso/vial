package vial

import (
    "context"
    "testing"

    gm "github.com/onsi/gomega"

    "github.com/daihasso/vial/responses"
)

func TestRouteControllerRespondsToMethod(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    handler := func(context.Context, *Transactor) responses.Data {
        return responses.Data{}
    }

    rch := RouteControllerHelper{
        route: Route{},
        methodCallers: map[RequestMethod]RouteControllerCaller{
            MethodGET: handler,
        },
    }

    responds := rch.RespondsToMethod(MethodGET)
    g.Expect(responds).To(gm.BeTrue())
}

func TestRouteControllerRespondsToMethodString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    handler := func(context.Context, *Transactor) responses.Data {
        return responses.Data{}
    }

    rch := RouteControllerHelper{
        route: Route{},
        methodCallers: map[RequestMethod]RouteControllerCaller{
            MethodGET: handler,
        },
    }

    responds := rch.RespondsToMethodString("get")
    g.Expect(responds).To(gm.BeTrue())
}
