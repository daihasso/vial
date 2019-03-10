package vial

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"

    "gopkg.in/yaml.v2"

    "daihasso.net/library/vial/responses"
    "daihasso.net/library/vial/neterr"
)

type SwaggerFormat int

const (
    _ SwaggerFormat = iota
    SwaggerYamlFormat
    SwaggerJsonFormat
)

var defaultSwaggerPath = "./swagger.yaml"

// defaultSwaggerController returns a dummy value of true.
type defaultSwaggerController struct {
    useJSON bool
    config  *Config
}

// Get will return a simple true for now.
func (self defaultSwaggerController) Get(
    ctx context.Context, transactor *Transactor,
) responses.Data {
    var bytes []byte
    var err error
    var contentType string

    path := self.getSwaggerPath()
    if self.useJSON {
        contentType = "application/json"
        bytes, err = readSwaggerDefinition(path)
    } else {
        contentType = "text/vnd.yaml"
        bytes, err = readSwaggerFile(path)
    }

    if os.IsNotExist(err) {
        return transactor.Abort(
            http.StatusNotFound, neterr.SwaggerNotFoundError,
        )
    } else if err != nil {
        transactor.Logger.Exception(
            err,
            "Error while reading swagger file.",
        )
        return transactor.Abort(
            http.StatusNotFound, neterr.SwaggerNotFoundError,
        )
    }

    return transactor.Respond(
        200,
        responses.Body(bytes),
        responses.AddHeader("Content-Type", contentType),
    )
}

// Get will return a simple true for now.
func (self *defaultSwaggerController) getSwaggerPath() string {
    path := defaultSwaggerPath
    if self.config.Swagger.Path != "" {
        path = self.config.Swagger.Path
    }

    return path
}

func readSwaggerDefinition(path string) ([]byte, error) {
    swaggerYAMLData, err := readSwaggerFile(path)
    if err != nil {
        return nil, err
    } else if swaggerYAMLData == nil {
        return nil, nil
    }

    var swagger interface{}

    err = yaml.Unmarshal(swaggerYAMLData, &swagger)
    if err != nil {
        return nil, err
    }

    // See: https://goo.gl/zI8Lph
    // NOTE: This doesn't preserve ordering. Should it?
    swagger = convert(swagger)

    jsonBytes, err := json.Marshal(swagger)
    if err != nil {
        return nil, err
    }

    return jsonBytes, nil
}

func readSwaggerFile(path string) ([]byte, error) {
    swaggerBytes, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    return swaggerBytes, nil
}

func convert(i interface{}) interface{} {
    switch x := i.(type) {
    case map[interface{}]interface{}:
        m2 := map[string]interface{}{}
        for k, v := range x {
            m2[k.(string)] = convert(v)
        }
        return m2
    case []interface{}:
        for i, v := range x {
            x[i] = convert(v)
        }
    }
    return i
}
