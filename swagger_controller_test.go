package vial

import (
    "os"
    "testing"
    "io/ioutil"
    "net/http"
    "net/http/httptest"

    gm "github.com/onsi/gomega"
    "github.com/daihasso/vial/responses"
)

var testSwaggerString = `
---
swagger: 2.0
info:
  description: A golang microservice.
  version: 0.0.1
  title: Test
host: test.com
basePath: /
paths:
  /health:
    get:
      description: Returns health.
      produces:
        - application/json
      responses:
        '200':
          description: Health status.
`[1:]


func TestSwaggerGet(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    swaggerFile, err := ioutil.TempFile("", "swagger.yaml")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(swaggerFile.Name())
    err = ioutil.WriteFile(swaggerFile.Name(), []byte(testSwaggerString), 0644)
    g.Expect(err).To(gm.BeNil())

    config := Config{
        Swagger: struct{
            Path string
        }{
            Path: swaggerFile.Name(),
        },
    }
    swaggerController := defaultSwaggerController{
        useJSON: false,
        config: &config,
    }

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    ctx, _ := handleSequenceId(req)
    req = req.WithContext(ctx)

    rr := httptest.NewRecorder()

    transactor, err := NewTransactor(
        req, rr, PathParams{}, &config, logger, responses.JSONEncoding,
    )
    g.Expect(err).To(gm.BeNil())

    response := swaggerController.Get(ctx, transactor)
    err = response.Write(rr)
    g.Expect(err).To(gm.BeNil())

    fileContents, _ := ioutil.ReadFile(swaggerFile.Name())
    t.Log(string(fileContents))

    g.Expect(rr.Body.String()).To(gm.BeEquivalentTo(testSwaggerString))
    g.Expect(rr.Header().Get("Content-Type")).To(
        gm.BeEquivalentTo("text/vnd.yaml"),
    )
}

func TestSwaggerGetJSON(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger := setupLogging(t, g)

    swaggerFile, err := ioutil.TempFile("", "swagger.yaml")
    g.Expect(err).To(gm.BeNil())
    defer os.Remove(swaggerFile.Name())
    err = ioutil.WriteFile(swaggerFile.Name(), []byte(testSwaggerString), 0644)
    g.Expect(err).To(gm.BeNil())

    config := Config{
        Swagger: struct{
            Path string
        }{
            Path: swaggerFile.Name(),
        },
    }
    swaggerController := defaultSwaggerController{
        useJSON: true,
        config: &config,
    }

    expectedBody := `{"basePath":"/","host":"test.com",` +
        `"info":{"description":"A golang microservice.",` +
        `"title":"Test","version":"0.0.1"},` +
        `"paths":{"/health":{"get":{"description":"Returns health.",` +
        `"produces":["application/json"],` +
        `"responses":{"200":{"description":"Health status."}}}}},"swagger":2}`

    req, err := http.NewRequest("GET", "/health", nil)
    g.Expect(err).To(gm.BeNil())

    ctx, _ := handleSequenceId(req)
    req = req.WithContext(ctx)

    rr := httptest.NewRecorder()

    transactor, err := NewTransactor(
        req, rr, PathParams{}, &config, logger, responses.JSONEncoding,
    )
    g.Expect(err).To(gm.BeNil())

    response := swaggerController.Get(ctx, transactor)
    err = response.Write(rr)
    g.Expect(err).To(gm.BeNil())

    fileContents, _ := ioutil.ReadFile(swaggerFile.Name())
    t.Log(string(fileContents))

    g.Expect(rr.Body.String()).To(gm.BeEquivalentTo(expectedBody))
    g.Expect(rr.Header().Get("Content-Type")).To(
        gm.BeEquivalentTo("application/json"),
    )
}
