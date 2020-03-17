package vial

import (
    "context"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    gm "github.com/onsi/gomega"

    "github.com/daihasso/vial/responses"
    "github.com/daihasso/vial/neterr"
)

type testRouteStruct struct {
    getMethod RouteFunction
}

func (self testRouteStruct) Get(transactor *Transactor) responses.Data {
    return self.getMethod(transactor)
}

type testRouteStructBad struct {}

func (self testRouteStructBad) Get(
    transactor *Transactor, bad string,
) responses.Data {
    return responses.Data{}
}

func TestAddPostRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("POST", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Post(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddGetRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("GET", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Get(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddPutRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("PUT", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Put(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddPatchRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("PATCH", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Patch(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddDeleteRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("DELETE", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Delete(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddHeadRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("HEAD", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Head(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddOptionsRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("OPTIONS", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Options(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestAddAllRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    server, err := NewServer(AddCustomLogger(logger))

    err = server.All(
        "/hello/<name>",
        func(transactor *Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }
            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    )
    g.Expect(err).To(gm.BeNil())

    for _, httpMethod := range validRouteActions {
        rr := httptest.NewRecorder()
        req, err := http.NewRequest(
            strings.ToUpper(httpMethod), "/hello/tester", nil,
        )
        g.Expect(err).To(gm.BeNil())

        server.muxer.ServeHTTP(rr, req)

        t.Log(rr.Result().Header)
        g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
        g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    }
}

func TestAddRouteWithContext(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("GET", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Get(
        "/hello/<name>",
        WithContext(func(
            ctx context.Context, transactor *Transactor,
        ) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    400,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }

            g.Expect(ctx).To(gm.Equal(transactor.Context()))

            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        }),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}


func TestAddRouteWithContextOnly(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `null`

    req, err := http.NewRequest("GET", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))

    err = server.Get(
        "/hello/<name>",
        ContextOnly(func(ctx context.Context) responses.Data {
            return responses.Respond(ctx, responses.JSONEncoding, 200)
        }),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
}
