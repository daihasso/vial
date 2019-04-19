package vial

import (
    "context"
    "io/ioutil"
    "math/rand"
    "net/http"
    "net/http/httptest"
    "os"
    "strconv"
    "strings"
    "testing"
    "time"

    "github.com/daihasso/slogging"
    gm "github.com/onsi/gomega"

    "github.com/daihasso/vial/responses"
    "github.com/daihasso/vial/neterr"
)

var testCertString = `-----BEGIN CERTIFICATE-----
MIIE+zCCAuOgAwIBAgIJALw8/UwIJNylMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCXRlc3QudGVzdDAeFw0xOTAyMjQxNzUyMjdaFw0yMDAyMjQxNzUyMjdaMBQx
EjAQBgNVBAMMCXRlc3QudGVzdDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBAMNmhi1hjz4Lg4Sb6i0tsofQT04ulSehWG60x9FVCqFzDRU7Bl2ILWHkEQAK
vJ+jfDO6Y/8Hxl1zY3M0o+ybSLfxP9HmcHifNNjFi09t3LbIDvU827ELFNfP8awE
fFkgJOvjLDV624fKoCnJdlPiEVzEYjErwbX05QC5DOrIjGddQK3AW7lEqxiMhoka
QCts03SPCOquJJCffbNbRGCoSNn/BzDNLjBNv8Z1oJmPGC8R7oCZn7QgWjY3OP+H
ZuCMRMtigRaGyPuay51oYiQm7ZeWTKE4qIsb/JAjD3BEaBMCuuj72DpXd+l8HYdq
K6f8Nojl7DcT5aoBcVb9lpxOcZyYu3YzS5cyQssToBNzDesquVVD1hC8ohKjnnIi
nvmwb427MmHztfV4h4xqRgKPhtJuk6+11sS0NJVkebVltavX/+zaiqvHcBfX9O0+
5ohBddYmQsVv3eXLdgwkFBCLL8DWhxWdUzECevPF0NqMpNySj3lRaAnEh0W9VCCv
IdtPW9Z3WJ3lkaqK06vDcvZJbJFX79GlnxCZZBxlaMA6TZsljLoKWOYQ0orB49c0
BFLfB96UJT2isGiFCFOw94UTMsm06cZycN+tsnFAZKsC9+JewHo56e5/6CLJYDMA
wYvfUadHMePieOxJOOsX+8/q5Z+KX9DFnV9i5ZXM0hUkwBBVAgMBAAGjUDBOMB0G
A1UdDgQWBBSlAY8X9x4pjAuX8FQ7ok2aQIdeozAfBgNVHSMEGDAWgBSlAY8X9x4p
jAuX8FQ7ok2aQIdeozAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQCt
GXTvOpf4nCShWXmQ1zXQJl8epybzQFMRBn6E0PaOcXpoTEzwNMjj+NTEbe1/38xb
f7bIYK+QP+UOlo/7/DxPfH2slGA5cnrWkfSUxLqNOX96TKWYpz9QmJFKgq5bOETy
mcE8UNVOylSujcsF6TdaDB1cxyiVtsd/IS2MSmi2dDSB/o9EoW0KdAIcofENvN18
JKR9YbMF4ZvjEO7zVJI0XKxvh1hii78wtki6lr/HskDzhow2BuL0SLTQJ3EcBcxW
e7HYYhkEbWMa3AnE3mQzSxXqzsshS2NUr0tvuCtjb1W1W6PuGKmqzxCgrVNaA3XM
LHE4NusEUtbJMur2lv0AEgBH3asjYz/EERcLPKRif6WK/WJFhnDKOSNoBagWlyCm
XzfQYWyLIaM5WRsESZWAm6qU3azh3ccJJuC7xgMrWcI2Kz8iaNhQMxwMP6cWsK3Z
/KDyOyccDEOHwQrT9HP8dsHD2GtcDFFaOTqZQcQhUIRuWx14+fh44AqGImFp3aAW
O4mCfZfQ0DpDgSbCU6u0gVhR0sWDKo9YcKy3Ah5sJsCQM/0zH3Gn/Gafb12pv93j
BaTLeOfnDUFts8qBezM8LnYLYy7Sb/7Si6e+X4vA1XqdftAeE+iuMgqECmLz+g9x
/XaKigpqwXdxUNmDnsk35HnAP64sp7DFA72FmlrW+g==
-----END CERTIFICATE-----
`

var testKeyString = `-----BEGIN PRIVATE KEY-----
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
A4mc4X5UAr2pMeu4wgzUWFF4KIIb3Mpm
-----END PRIVATE KEY-----
`

func setupLogging(t *testing.T, g *gm.GomegaWithT) *logging.Logger {
    err := logging.GetRootLogger().SetLogLevel(logging.DEBUG)
    g.Expect(err).To(gm.BeNil())
    logging.GetRootLogger().SetWriters(testWriter{t})

    identifier := "test" + strconv.Itoa(rand.Int()) //#nosec G404
    logger, err := logging.NewLogger(identifier)
    g.Expect(err).To(gm.BeNil())

    return logger
}

type testWriter struct {
    t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
    tw.t.Log(string(p))
    return len(p), nil
}

func TestServerBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"healthy":true}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestServerBasicEncrypted(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"healthy":true}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    testCertData := strings.NewReader(testCertString)
    testKeyData := strings.NewReader(testKeyString)
    rr := httptest.NewRecorder()
    server, err := NewServerDefault(
        AddCustomLogger(logger),
        AddEncryption(testCertData, testKeyData),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
    g.Expect(rr.Header().Get("Strict-Transport-Security")).To(gm.Equal(
        "max-age=63072000; includeSubDomains",
    ))
}

func TestServerFuncMultiRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("GET", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())
    req2, err := http.NewRequest("POST", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/hello/<name>",
        FuncHandlerMulti(
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
            "get",
            "post",
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())

    rr = httptest.NewRecorder()

    server.muxer.ServeHTTP(rr, req2)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestServerSimpleRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"hello":"tester"}`

    req, err := http.NewRequest("GET", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/hello/<name>",
        FuncHandler(
            "get",
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
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}

func TestServerDefaultOptions(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("OPTIONS", "/hello/tester", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/hello/<name>",
        FuncHandler(
            "get",
            func(transactor *Transactor) responses.Data {
                return transactor.Respond(200)
            },
        ),
    )
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/hello/<name>",
        FuncHandler(
            "post",
            func(transactor *Transactor) responses.Data {
                return transactor.Respond(
                    200,
                )
            },
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Header().Get("Access-Control-Allow-Methods")).To(
        gm.Equal("OPTIONS, GET, POST"),
    )
}

func TestServerBadMethod(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("PUT", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    t.Log(rr.Body.String())
    g.Expect(rr.Code).To(gm.Equal(http.StatusMethodNotAllowed))
    g.Expect(rr.Header().Get("Allow")).To(
        gm.Equal("OPTIONS, GET"),
    )
}

func TestServerBasicGoStart(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    server, err := NewServerDefault(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    ch := server.GoStart()

    for resp := range ch {
        t.Log("S resp", resp)
        switch resp.Type {
            case ServerStartChannelResponse:
            time.Sleep(1 * time.Millisecond)
            err := server.Stop()
            g.Expect(err).To(gm.BeNil())
            continue
            case ServerShutdownChannelResponse:
            return
            default:
            t.Fatal("Received unexpected response from server:", resp)
        }
    }
}

func TestServerEncryptedGoStart(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    testCertData := strings.NewReader(testCertString)
    testKeyData := strings.NewReader(testKeyString)
    server, err := NewServerDefault(
        AddCustomLogger(logger),
        AddEncryption(testCertData, testKeyData),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    ch := server.GoStart()

    for resp := range ch {
        t.Log("S resp", resp)
        switch resp.Type {
            case ServerStartChannelResponse:
            time.Sleep(1 * time.Millisecond)
            err := server.Stop()
            g.Expect(err).To(gm.BeNil())
            continue
            case ServerShutdownChannelResponse:
            return
            default:
            t.Fatal("Received unexpected response from server:", resp)
        }
    }
}

func TestServerInvalidRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/test/route",
        FuncHandler(
            "get",
            func(foo string, bar int) {},
        ),
    )
    g.Expect(err).To(gm.HaveOccurred())
}

func TestServerMissingRoute(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/foo", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusNotFound))
    // TODO: Need to make sequence ids return always. This requires somehow
    // hijacking the default 404 handler. It will take some work. It might
    // even require a custom handler.
    // g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
}


func TestServerPreActionModifyContext(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test/middleware", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(
        AddPreActionMiddleware(
            func(ctx context.Context, transactor *Transactor) (
                *responses.Data, *context.Context, error,
            ) {
                newCtx := context.WithValue(ctx, "foo", "bar")
                return nil, &newCtx, nil
            }),
        AddCustomLogger(logger),
    )
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "test/middleware",
        FuncHandler(
            "get",
            func(ctx context.Context, transactor *Transactor) responses.Data {
                return transactor.Respond(
                    200,
                    responses.Body(map[string]string{
                        "context_value": ctx.Value("foo").(string),
                    }),
                )
            },
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(`{"context_value":"bar"}`))
}

func TestServerPreActionHijackReturn(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(AddPreActionMiddleware(
        func(_ context.Context, transactor *Transactor) (
            *responses.Data, *context.Context, error,
        ) {
            return &responses.Data{
                Body: []byte(`HIJACKED!`),
                StatusCode: http.StatusOK,
            }, nil, nil
        }),
        AddCustomLogger(logger),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(`HIJACKED!`))
}

func TestServerPostActionHijackReturn(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    expectedBody := `{"healthy":true}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(AddPostActionMiddleware(
        func(
            _ context.Context,
            transactor *Transactor,
            existingResp responses.Data,
        ) (*responses.Data, error) {
            return &responses.Data{
                Headers: existingResp.Headers,
                Body: existingResp.Body,
                StatusCode: http.StatusAccepted,
            }, nil
        }),
        AddCustomLogger(logger),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Code).To(gm.Equal(http.StatusAccepted))
}

func TestServerBasicEncryptedFilePaths(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    tempCertFile, err := ioutil.TempFile("", "test-server-encryption.cert")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(tempCertFile.Name())
    err = ioutil.WriteFile(tempCertFile.Name(), []byte(testCertString), 0644)
    g.Expect(err).To(gm.BeNil())

    tempKeyFile, err := ioutil.TempFile("", "test-server-encryption.key")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(tempKeyFile.Name())
    err = ioutil.WriteFile(tempKeyFile.Name(), []byte(testKeyString), 0644)
    g.Expect(err).To(gm.BeNil())

    expectedBody := `{"healthy":true}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(
        AddCustomLogger(logger),
        AddEncryptionFilePaths(tempCertFile.Name(), tempKeyFile.Name()),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
    g.Expect(rr.Header().Get("Strict-Transport-Security")).To(gm.Equal(
        "max-age=63072000; includeSubDomains",
    ))
}

func TestServerBasicEncryptedFilePathsViaConfig(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    tempCertFile, err := ioutil.TempFile("", "test-server-encryption.cert")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(tempCertFile.Name())
    err = ioutil.WriteFile(tempCertFile.Name(), []byte(testCertString), 0644)
    g.Expect(err).To(gm.BeNil())

    tempKeyFile, err := ioutil.TempFile("", "test-server-encryption.key")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(tempKeyFile.Name())
    err = ioutil.WriteFile(tempKeyFile.Name(), []byte(testKeyString), 0644)
    g.Expect(err).To(gm.BeNil())

    expectedBody := `{"healthy":true}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    config := newConfig()
    config.Tls.CertPath = tempCertFile.Name()
    config.Tls.KeyPath = tempKeyFile.Name()
    config.Tls.Enabled = true

    rr := httptest.NewRecorder()
    server, err := NewServerDefault(
        AddCustomLogger(logger),
        AddConfig(config),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    t.Log(rr.Result().Header)
    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Sequence-Id")).ToNot(gm.BeEmpty())
    g.Expect(rr.Header().Get("Strict-Transport-Security")).To(gm.Equal(
        "max-age=63072000; includeSubDomains",
    ))
}

func TestServerCustomConfigFile(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    configYaml := `
vial:
  port: 8081
`
    tempConfigFile, err := ioutil.TempFile("", "config.yaml")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(tempConfigFile.Name())
    err = ioutil.WriteFile(tempConfigFile.Name(), []byte(configYaml), 0644)
    g.Expect(err).To(gm.BeNil())

    server, err := NewServerDefault(
        AddCustomLogger(logger), AddConfigFromFile(tempConfigFile.Name()),
    )
    g.Expect(err).To(gm.BeNil())

    g.Expect(server.GetConfig().Port).To(gm.Equal(8081))

    // Host should be default value.
    g.Expect(server.GetConfig().Host).To(gm.Equal("127.0.0.1"))
}

func TestServerLoggerContext(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/test",
        FuncHandler(
            "get",
            func(ctx context.Context, transactor *Transactor) responses.Data {
                serverLoggerIn := ctx.Value(ServerLoggerContextKey)
                serverLogger := serverLoggerIn.(*logging.Logger)

                g.Expect(serverLogger == serverLogger)

                t.Log("Logger name:", serverLogger.Identifier())

                return transactor.Respond(
                    200,
                )
            },
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
}

func TestTransactorFromContext(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())

    err = server.AddController(
        "/test",
        FuncHandler(
            "get",
            func(ctx context.Context, transactor *Transactor) responses.Data {
                ctxTransactorIn := ctx.Value(TransactorContextKey)
                ctxTransactor := ctxTransactorIn.(*Transactor)

                g.Expect(ctxTransactor == ctxTransactor)

                return transactor.Respond(200)
            },
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
}

func TestUrlFor(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test2", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())
    handler := func(transactor *Transactor) responses.Data {
        return transactor.Respond(200)
    }

    handler2 := func(transactor *Transactor) responses.Data {
        return transactor.Respond(
            200,
            responses.Body(transactor.UrlFor(handler)),
        )
    }

    err = server.AddController(
        "/test",
        FuncHandler(
            "get",
            handler,
        ),
    )
    g.Expect(err).To(gm.BeNil())
    err = server.AddController(
        "/test2",
        FuncHandler(
            "get",
            handler2,
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal("/test"))
}

func TestUrlForVariableSub(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test2", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())
    handler := func(transactor *Transactor) responses.Data {
        return transactor.Respond(200)
    }

    handler2 := func(transactor *Transactor) responses.Data {
        return transactor.Respond(
            200,
            responses.Body(transactor.UrlFor(handler, UrlParamValues{
                "foo": 5,
            })),
        )
    }

    err = server.AddController(
        "/test/<int:foo>",
        FuncHandler(
            "get",
            handler,
        ),
    )
    g.Expect(err).To(gm.BeNil())
    err = server.AddController(
        "/test2",
        FuncHandler(
            "get",
            handler2,
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal("/test/5"))
}

func TestUrlForVariableSubByIndex(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test2", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())
    handler := func(transactor *Transactor) responses.Data {
        return transactor.Respond(200)
    }

    handler2 := func(transactor *Transactor) responses.Data {
        return transactor.Respond(
            200,
            responses.Body(transactor.UrlFor(handler, 5)),
        )
    }

    err = server.AddController(
        "/test/<int:foo>",
        FuncHandler(
            "get",
            handler,
        ),
    )
    g.Expect(err).To(gm.BeNil())
    err = server.AddController(
        "/test2",
        FuncHandler(
            "get",
            handler2,
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal("/test/5"))
}

func TestUrlForVariableSubWrongKey(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    req, err := http.NewRequest("GET", "/test2", nil)
    g.Expect(err).To(gm.BeNil())

    rr := httptest.NewRecorder()
    server, err := NewServer(AddCustomLogger(logger))
    g.Expect(err).To(gm.BeNil())
    handler := func(transactor *Transactor) responses.Data {
        return transactor.Respond(200)
    }

    handler2 := func(transactor *Transactor) responses.Data {
        return transactor.Respond(
            200,
            responses.Body(transactor.UrlFor(handler, UrlParamValues{
                "bar": 5,
            })),
        )
    }

    err = server.AddController(
        "/test/<int:foo>",
        FuncHandler(
            "get",
            handler,
        ),
    )
    g.Expect(err).To(gm.BeNil())
    err = server.AddController(
        "/test2",
        FuncHandler(
            "get",
            handler2,
        ),
    )
    g.Expect(err).To(gm.BeNil())

    server.muxer.ServeHTTP(rr, req)

    g.Expect(rr.Code).To(gm.Equal(http.StatusOK))
    g.Expect(rr.Body.String()).To(gm.Equal("/test/<int:foo>"))
}
