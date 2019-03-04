package vial

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"

    "github.com/daihasso/slogging"
    "github.com/pkg/errors"
    "github.com/google/uuid"

    "daihasso.net/library/vial/responses"
    "daihasso.net/library/vial/neterr"
)

// Transactor is a helper for handling request and response logic.
type Transactor struct {
    Builder *responses.Builder
    Config *Config
    Logger logging.Logger
    Request *Request

    context context.Context
}

func marshalItem(item interface{}) (string, error) {
    var outputString string
    var marshalledOutput []byte
    var err error

    marshalledOutput, err = json.Marshal(item)
    if err != nil {
        return "", errors.Wrap(err, "Error marshalling item to JSON")
    }

    outputString = string(marshalledOutput)

    if outputString == "{}" {
        outputString = ""
    }

    return outputString, nil
}

// AddHeader marshales content and adds it as a header key.
func (t *Transactor) AddHeader(key string, content interface{}) error {
    var output string
    if contentString, ok := content.(string); ok {
        output = contentString
    } else {
        marshalledContent, err := marshalItem(content)
        if err != nil {
            return errors.Wrap(err, "Error marshalling header to JSON")
        }
        output = string(marshalledContent)
    }
    t.Builder.AddHeader(key, output)

    return nil
}

// SetHeader marshales content and sets key's content to it.
func (t *Transactor) SetHeader(key string, content interface{}) error {
    var output string
    if contentString, ok := content.(string); ok {
        output = contentString
    } else {
        marshalledContent, err := json.Marshal(content)
        if err != nil {
            return errors.Wrap(err, "Error marshalling header to JSON")
        }
        output = string(marshalledContent)
    }

    t.Builder.SetHeader(key, output)

    return nil
}

/*
/* TODO: Re-enable after JWT library is done.
// AddJWT adds JWT as authorization.
func (t *Transactor) AddJWT(jwt *JWT) {
    jwtString, err := jwt.Finalize()
    if err != nil {
        panic(err)
    }
    t.Builder.SetHeader(
        "Authorization",
        fmt.Sprintf("Bearer %s", jwtString),
    )
}
*/

// RequestBodyString return the request body as a string.
func (i *Transactor) RequestBodyString() (string, error) {
    // TODO: Check performance on this, maybe read and store the data as a
    // string and act upon that string instead of duplicating the read every
    // time.
    var buf bytes.Buffer

    req := i.Request
    bodyReader := req.Body
    bodyTeeReader := io.TeeReader(bodyReader, &buf)

    bodyString, err := ioutil.ReadAll(bodyTeeReader)
    if err != nil {
        return "", err
    }

    return string(bodyString), nil
}

// CopyHeadersFromResponse will copy the headers from a client response into
// the server response.
func (t *Transactor) CopyHeadersFromResponse(
    response *http.Response,
) {
    for k, v := range response.Header {
        if len(v) > 1 {
            t.Builder.AddHeader(k, v[0], v[1:]...)
        } else {
            t.Builder.AddHeader(k, v[0])
        }
    }
}

// Respond is a shortcut for finishing and responding to a request.
func (self Transactor) Respond(
    statusCode int, attributes ...responses.AdditionalAttribute,
) responses.Data {
    self.Builder.SetStatus(statusCode)
    return self.Builder.Finish(attributes...)
}

// Abort is a shortcut for responding to a request with a failure.
func (self Transactor) Abort(
    statusCode int,
    codedErr neterr.CodedError,
    otherErrors ...neterr.CodedError,
) responses.Data {
    err := self.SetHeader("Content-Type", responses.JSONContentType)
    if err != nil {
        return responses.ErrorResponse(err)
    }
    return self.Builder.Abort(statusCode, codedErr, otherErrors...)
}

// SequenceId grabs the current Sequence ID from the context.
func (i Transactor) SequenceId() *uuid.UUID {
    // NOTE: This method expects the SequenceId to be in the context by now. It
    //       may end up returning nil if something goes awry.
    sequenceId, _ := ContextSequenceId(i.context)
    return sequenceId
}

// NewTransactor will generate a new transactor for request to a controller.
func NewTransactor(
    request *http.Request,
    responseWriter http.ResponseWriter,
    variables map[string]string,
    config *Config,
    logger logging.Logger,
    encodingType responses.EncodingType,
    options ...TransactorOption,
) (*Transactor, error) {
    existingContext := request.Context()
    sequenceId, _ := ContextSequenceId(existingContext)
    builder, err := responses.NewBuilder(
        existingContext,
        encodingType,
        responses.AddHeader(SequenceIdHeader, sequenceId.String()),
    )
    if err != nil {
        return nil, errors.Wrap(
            err, "Error while initializing response Builder",
        )
    }

    newRequest := NewServerRequest(
        request,
        variables,
    )
    transactor := &Transactor{
        Builder: builder,
        Config: config,
        Request: newRequest,
        Logger: nil,

        context: existingContext,
    }

    loggerWithSequenceId := logging.WithExtras(
        logger, logging.AddExtraAttributeFuncs(
            map[string]func() (interface{}, error){
                "sequence_id": func() (interface{}, error) {
                    return transactor.SequenceId().String(), nil
                },
            },
        ),
    )
    transactor.Logger = loggerWithSequenceId

    for _, option := range options {
        err := option(transactor)
        if err != nil {
            return nil, errors.Wrap(
                err, "Error while setting options on Transactor",
            )
        }
    }

    return transactor, nil
}
