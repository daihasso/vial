package vial

import (
    "context"
    "fmt"
    "os"
    "reflect"
    "sort"
    "strings"
    "strconv"
    "log"
    "net/http"
    "crypto/tls"
    "io"
    "io/ioutil"

    "github.com/pkg/errors"
    "github.com/daihasso/slogging"
    "github.com/daihasso/peechee"
    "github.com/google/uuid"

    "github.com/daihasso/vial/responses"
    "github.com/daihasso/vial/neterr"
)

// Server is a specialized server for microservices.
type Server struct {
    PathReader *peechee.PathReader
    Logger *logging.Logger

    config *Config
    muxer *http.ServeMux
    urlForMap map[reflect.Value][]string
    pathRouteControllerHelpers map[string][]*RouteControllerHelper
    preActionMiddleware []PreMiddleWare
    postActionMiddleware []PostMiddleWare
    internalServer *http.Server
    defaultEncoding responses.EncodingType
    encryptionEnabled bool
}

func setupTls(tlsCertData, tlsKeyData io.Reader) (*tls.Config, error) {
    certData, err := ioutil.ReadAll(tlsCertData)
    if err != nil {
        return nil, errors.Wrap(err, "Couldn't get certificate data")
    }
    keyData, err := ioutil.ReadAll(tlsKeyData)
    if err != nil {
        return nil, errors.Wrap(err, "Couldn't get private key data")
    }

    cert, err := tls.X509KeyPair(certData, keyData)
    if err != nil {
        return nil, errors.Wrap(
            err, "Error while adding cert and key to keypair",
        )
    }

    return &tls.Config{
        MinVersion: tls.VersionTLS12,
        CurvePreferences: []tls.CurveID{
            tls.CurveP521,
            tls.CurveP384,
            tls.CurveP256,
        },
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
        Certificates: []tls.Certificate{cert},
    }, nil
}

func createGoServer(
    host string,
    port int,
    muxer *http.ServeMux,
    tlsConfig *tls.Config,
    logger *logging.Logger,
) *http.Server {
    handleProfiling(muxer)
    return &http.Server{
        Addr: fmt.Sprintf("%s:%s", host, strconv.Itoa(port)),
        Handler: muxer,
        TLSConfig: tlsConfig,
        TLSNextProto: make(
            map[string]func(*http.Server, *tls.Conn, http.Handler),
        ),
        ErrorLog: log.New(
            logging.NewPseudoWriter(logging.ERROR, logger), "", 0,
        ),
    }
}

// AddController adds a new controller to the server at the route specified. It
// must match the expected format that RouteControllerIsValid checks for.
// It must be one of:
//    A func or a struct with a method defined in validRouteControllerFields
//    (Post, Get, Put, Patch, Delete, Head and/or Options) Options method that
//    matches either of the signatures defined in validRouteControllerTypes:
//        func(context.Context, *Transactor) responses.Data
//    or
//        func(*Transactor) responses.Data
func (s *Server) AddController(
    path string,
    rc RouteController,
    otherRCs ...RouteController,
) error {
    if errs := RouteControllerIsValid(rc); len(errs) != 0 {
        return errors.Errorf(
            "RouteController '%T' provided to AddController is not valid:\n%s",
            rc,
            strings.Join(errs, "\nand\n"),
        )
    }

    route, err := ParseRoute(path)
    if err != nil {
        return errors.Wrap(err, "Error while parsing route provided")
    }
    allRouteControllers := append([]RouteController{rc}, otherRCs...)
    methodCallers, urlForMap := MethodsForRouteController(
        path, allRouteControllers...,
    )
    routeControllerHelper := RouteControllerHelper{
        route: route,
        methodCallers: methodCallers,
    }

    for k, v := range urlForMap {
        if existing, ok := s.urlForMap[k]; ok {
            s.urlForMap[k] = append(existing, v)
        } else {
            s.urlForMap[k] = []string{path}
        }
    }

    if _, ok := s.pathRouteControllerHelpers[route.Base]; !ok {
        s.pathRouteControllerHelpers[route.Base] = make(
            []*RouteControllerHelper,
            0,
        )

        s.muxer.HandleFunc(
            route.Base,
            s.defaultMultiRouteControllerWrapper(route.Base),
        )
    }

    s.pathRouteControllerHelpers[route.Base] = append(
        s.pathRouteControllerHelpers[route.Base],
        &routeControllerHelper,
    )

    return nil
}

// GetConfig retrieves the config by value.
func (self Server) GetConfig() Config {
    return *self.config
}


// Start starts the server. This call is blocking.
func (s *Server) Start() error {
    return s.startUp()
}

// GoStart starts the server. This call is non-blocking.
func (self *Server) GoStart() chan ServerChannelResponse {
    outCh := make(chan ServerChannelResponse)
    go self.goStartUp(outCh)
    return outCh
}

// Stop stops the server.
func (s *Server) Stop() error {
    s.Logger.Info("Stopping server...")
    err := s.internalServer.Close()
    if err != nil {
        return errors.Wrap(err, "Error while stopping server")
    }

    return nil
}

func (s Server) startInternalServer() error {
    s.Logger.Info("Starting server...", logging.Extras{
            "host": s.config.Host,
            "port": s.config.Port,
            "using_encryption": s.encryptionEnabled,
    })
    if s.encryptionEnabled {
        return s.internalServer.ListenAndServeTLS("", "")
    }

    return s.internalServer.ListenAndServe()
}

func (s Server) goStartUp(outCh chan ServerChannelResponse) {
    outCh <- ServerChannelResponse{
        Type: ServerStartChannelResponse,
        Error: nil,
    }

    err := s.startInternalServer()

    if err != http.ErrServerClosed {
        outCh <- ServerChannelResponse{
            Type: UnknownErrorChannelResponse,
            Error: err,
        }
    } else {
        s.Logger.Info("Server shutdown.")
        outCh <- ServerChannelResponse{
            Type: ServerShutdownChannelResponse,
            Error: nil,
        }
    }
}

func (s *Server) startUp() error {
    err := s.startInternalServer()
    if err != http.ErrServerClosed {
        s.Logger.Exception(
            err,
            "Error while trying to run server.",
            logging.Extras{
                "tls_config": s.internalServer.TLSConfig,
            },
        )
        return errors.Wrap(err, "Error while running server")
    }

    s.Logger.Info("Server shutdown.")
    return nil
}

func (self Server) getMatchingControllers(
    requestPath string,
    basePath string,
) []*RouteControllerHelper {
    var methodHelpers []*RouteControllerHelper
    for _, rch := range self.pathRouteControllerHelpers[basePath] {
        if rch.route.Matches(requestPath) {
            methodHelpers = append(methodHelpers, rch)
        }
    }

    return methodHelpers
}

func (self Server) respondToMethod(
    w http.ResponseWriter, r *http.Request, rchs []*RouteControllerHelper,
) responses.Data {
    reqMethod := RequestMethodFromString(r.Method)

    var rcc RouteControllerCaller
    rch := rchs[0]
    rchSet := false
    for _, nextRch := range rchs {
        if nextRcc, ok := rch.ControllerFuncForMethod(reqMethod); ok {
            rch = nextRch
            rcc = nextRcc
            rchSet = true
            break
        }
    }

    if !rchSet {
        if reqMethod == MethodOPTIONS {
            rcc = func(
                _ context.Context, transactor *Transactor,
            ) responses.Data {
                return DefaultOptions(self, transactor, rchs)
            }
        } else {
            //
            // NOTE: A lot of the logic here of iterating through all the
            //       helpers could probably be optimized by doing this at
            //       AddController time and storing a reverse-map of route
            //       matchers to methods or something like that.
            //
            //       I'm not going to optimize it yet because it seems
            //       premature but this may be worth revisiting.
            //
            self.Logger.Warn(
                fmt.Sprintf(
                    "You have not set up your %s method for this route.",
                    reqMethod.String(),
                ),
            )

            var methodStrings []string
            methodMap := map[RequestMethod]bool{
                // OPTIONS is always supported.
                MethodOPTIONS: true,
            }
            for _, nextRch := range rchs {
                for _, method := range nextRch.AllMethods() {
                    if _, ok := methodMap[method]; !ok {
                        methodStrings = append(methodStrings, method.String())
                        methodMap[method] = true
                    }
                }
            }

            sort.Strings(methodStrings)
            methodStrings = append(
                []string{MethodOPTIONS.String()}, methodStrings...,
            )
            allowedMethods := strings.Join(methodStrings, ", ")
            // TODO: Maybe check the error here. Can it actually occur?
            resp, _ := responses.NewBuilder(
                r.Context(),
                self.defaultEncoding,
                responses.Headers(map[string][]string{
                    "Allow": []string{allowedMethods},
                }),
            )
            return resp.Abort(
                http.StatusMethodNotAllowed,
                neterr.MethodNotAllowedError,
            )
        }
    }

    pathVariables, err := rch.route.PathParams(r.URL.Path)
    if err != nil {
        self.Logger.Exception(err, "Failed to parse path params properly.")
        sequenceId, err := ContextSequenceId(r.Context())
        if err != nil {
            return responses.ErrorResponse(err)
        }
        builder, err := responses.NewBuilder(
            r.Context(),
            self.defaultEncoding,
            responses.AddHeader(
                SequenceIdHeader,
                sequenceId.String(),
            ),
        )
        if err != nil {
            return responses.ErrorResponse(err)
        }
        return builder.Abort(
            http.StatusNotFound,
            neterr.RouteNotSetupError,
        )
    }
    transactor, err := NewTransactor(
        r, w, pathVariables, self.config, self.Logger, self.defaultEncoding,
    )
    if err != nil {
        self.Logger.Exception(err, "Error while creating Transactor.")
        return responses.ErrorResponse(err)
    }

    for _, middleware := range self.preActionMiddleware {
        data, newCtx, err := middleware(r.Context(), transactor)
        if err != nil {
            self.Logger.Exception(err, "Error in pre-action middleware.")
            return responses.ErrorResponse(err)
        }
        if data != nil {
            // If we have data from our middleware return early with it.
            return *data
        }
        if newCtx != nil {
            transactor.ChangeContext(*newCtx)
            r = &transactor.Request.Request
        }
    }

    self.Logger.Debug("Handling HTTP request.", logging.Extras{
        "path": transactor.Request.URL.Path,
        "method": transactor.Request.Method,
        "sequence_id": transactor.SequenceId(),
    })

    response := rcc(transactor.Context(), transactor)

    for _, middleware := range self.postActionMiddleware {
        data, err := middleware(transactor.Context(), transactor, response)
        if err != nil {
            self.Logger.Exception(err, "Error in pre-action middleware.")
            return responses.ErrorResponse(err)
        }
        if data != nil {
            // If we have data to return early with it.
            return *data
        }
    }

    transactor.Logger.Close()

    return response
}

func handleSequenceId(r *http.Request) (context.Context, string) {
    sequenceIdString := r.Header.Get(SequenceIdHeader)
    sequenceId, parseErr := uuid.Parse(sequenceIdString)
    if parseErr != nil {
        var err error
        sequenceId, err = uuid.NewRandom()
        if err != nil {
            panic(errors.Wrap(
                err, "Error while attempting to generate SequenceId",
            ))
        }
    }

    return ContextWithSequenceId(r.Context(), sequenceId), sequenceId.String()
}

func handleRequestId(ctx context.Context) context.Context {
    requestId, err := uuid.NewRandom()
    if err != nil {
        panic(errors.Wrap(
            err, "Failed to generated request id",
        ))
    }

    return context.WithValue(
        ctx, RequestIdContextKey, requestId.String(),
    )
}

type requestHandlerFunc func(
    w http.ResponseWriter, r *http.Request,
) responses.Data

func responseProcessor(
    handlerFunc requestHandlerFunc, server *Server,
) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var sequenceId string
        var ctx context.Context
        defer func() {
            if rawErr := recover(); rawErr != nil {
                newErr := errors.New(fmt.Sprintf("%+v", rawErr))
                if err, ok := rawErr.(error); ok {
                    newErr = errors.WithStack(err)
                }
                server.Logger.Exception(
                    newErr, "Panic while handling controller.",
                )
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "Internal Server Error")
            }
        }()

        // Add a reference to ourself to the context.
        r = r.WithContext(context.WithValue(
            r.Context(), ServerContextKey, server,
        ))

        // Add our sequence id, request id & server logger to the context.
        ctx, _ = handleSequenceId(r)
        ctx = handleRequestId(ctx)
        ctx = context.WithValue(ctx, ServerLoggerContextKey, server.Logger)
        r = r.WithContext(ctx)

        responseData := handlerFunc(w, r)
        if unexpectedErr := responseData.Error(); unexpectedErr != nil {
            if unexpectedErr.Error() != "" {
                server.Logger.Exception(
                    unexpectedErr,
                    "Unexpected error in controller route.",
                    logging.Extras{
                        "sequence_id": sequenceId,
                    },
                )
            }
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "Internal Server Error")
            return
        }

        err := responseData.Write(w)
        panic(err)
    }
}

func (s *Server) defaultMultiRouteControllerWrapper(
    basePath string,
) func(http.ResponseWriter, *http.Request) {
    requestFuncMatcher := func(
        w http.ResponseWriter,
        r *http.Request,
    ) responses.Data {
        // Add our sequence id to the context.
        sequenceId, err := ContextSequenceId(r.Context())
        if err != nil {
            s.Logger.Warn(
                "Unable to extract SequenceId from context.",
                logging.Extras{
                    "error": err,
                },
            )
        }
        s.Logger.Info("HTTP request made.", logging.Extras{
            "path": r.URL.Path,
            "requestor": r.RemoteAddr,
            "method": r.Method,
            "sequence_id": sequenceId,
        })

        matchingRouteControllerHelpers := s.getMatchingControllers(
            r.URL.Path, basePath,
        )
        if matchingRouteControllerHelpers == nil {
            builder, err := responses.NewBuilder(
                r.Context(),
                s.defaultEncoding,
                responses.AddHeader(
                    SequenceIdHeader,
                    sequenceId.String(),
                ),
            )
            if err != nil {
                return responses.ErrorResponse(err)
            }

            return builder.Abort(
                http.StatusNotFound,
                neterr.RouteNotSetupError,
            )
        }

        return s.respondToMethod(w, r, matchingRouteControllerHelpers)
    }

    return responseProcessor(requestFuncMatcher, s)
}

func (self Server) UrlFor(handler interface{}) string {
    if urls := self.UrlsFor(handler); urls != nil {
        return urls[0]
    }

    return ""
}

func (self Server) UrlsFor(handler interface{}) []string {
    handlerVal := reflect.ValueOf(handler)
    for handlerVal.Kind() == reflect.Ptr {
        handlerVal = handlerVal.Elem()
    }
    if urls, ok := self.urlForMap[handlerVal]; ok {
        return urls
    }

    return nil
}

// NewServer creates a bare-bones server.
func NewServer(options ...ServerOption) (*Server, error) {
    return newServer(options)
}

// NewServerDefault creates a new server with (mostly) sensible defaults.
func NewServerDefault(options ...ServerOption) (*Server, error) {
    allOptions := append(
        options,
        []ServerOption{
            AddDefaultSwaggerRoute(SwaggerYamlFormat, SwaggerJsonFormat),
            AddDefaultHealthRoute(),
        }...,
    )
    server, err :=  newServer(allOptions)
    if err != nil {
        return nil, err
    }

    return server, nil
}

func newServer(options []ServerOption) (*Server, error) {
    svOpts := newServerOptions()

    for _, option := range options {
        err := option(svOpts)
        if err != nil {
            return nil, errors.WithStack(err)
        }
    }

    config := svOpts.config
    if config == nil {
        config = newConfig()
    }


    useEncryption := svOpts.useEncryption || config.Tls.Enabled

    logger := svOpts.logger
    if logger == nil {
        var err error
        logger, err = logging.NewLogger(
            "vial.server.logger",
            logging.WithFormat(logging.JSON),
            logging.WithLogWriters(os.Stdout),
            logging.WithLogLevel(logging.INFO),
        )
        if err != nil {
            return nil, errors.Wrap(
                err, "Error while trying to create logger for server",
            )
        }
    }

    preActionMiddleware := svOpts.preActionMiddleware
    postActionMiddleware := svOpts.postActionMiddleware

    var tlsCfg *tls.Config
    if useEncryption {
        preActionMiddleware = append(
            preActionMiddleware, DefaultEncryptionHeadersMiddleware(),
        )

        if svOpts.tlsCertData == nil {
            configCertPath := config.Tls.CertPath
            if configCertPath == "" {
                return nil, errors.New(
                    "Encryption enabled but certificate path or data not" +
                        " provided",
                )
            }

            certReader, err := svOpts.pathReader.Read(configCertPath)
            if err != nil {
                return nil, errors.Wrap(
                    err, "Error while reading cert at path provided in config",
                )
            }

            svOpts.tlsCertData = certReader
        }

        if svOpts.tlsKeyData == nil {
            configKeyPath := config.Tls.KeyPath
            if configKeyPath == "" {
                return nil, errors.New(
                    "Encryption enabled but key path or data not provided",
                )
            }

            keyReader, err := svOpts.pathReader.Read(configKeyPath)
            if err != nil {
                return nil, errors.Wrap(
                    err, "Error while reading key at path provided in config",
                )
            }

            svOpts.tlsKeyData = keyReader
        }

        var err error
        tlsCfg, err = setupTls(svOpts.tlsCertData, svOpts.tlsKeyData)
        if err != nil {
            return nil, errors.Wrap(err, "Error while setting up TLS config")
        }
    }


    muxer := http.NewServeMux()

    goServer := createGoServer(
        config.Host,
        config.Port,
        muxer,
        tlsCfg,
        logger,
    )

    defaultEncoding := svOpts.defaultEncoding
    if defaultEncoding == responses.UnsetEncoding {
        // NOTE: We fallback to JSON if it's not explicitly set, this may be a
        //       little opinionated but JSON is the first-class citizen here.
        defaultEncoding = responses.JSONEncoding
    }


    server := &Server{
        PathReader: svOpts.pathReader,
        Logger: logger,

        config: config,
        muxer: muxer,
        urlForMap: make(map[reflect.Value][]string),
        pathRouteControllerHelpers: make(map[string][]*RouteControllerHelper),
        preActionMiddleware: preActionMiddleware,
        postActionMiddleware: postActionMiddleware,
        internalServer: goServer,
        defaultEncoding: defaultEncoding,
        encryptionEnabled: useEncryption,
    }

    for _, mod := range svOpts.serverMods {
        err := mod(server)
        if err != nil {
            return nil, errors.Wrap(
                err, "Error while running server modification option",
            )
        }
    }


    return server, nil
}
