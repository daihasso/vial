package vial

import (
    "context"
    "fmt"
    "sort"
    "strings"
    "strconv"
    "net/http"
    "crypto/tls"
    "io"
    "io/ioutil"

    "github.com/pkg/errors"
    "github.com/daihasso/slogging"
    "github.com/daihasso/peechee"
    "github.com/google/uuid"

    "daihasso.net/library/vial/responses"
    "daihasso.net/library/vial/neterr"
)

// Server is a specialized server for microservices.
type Server struct {
    PathReader *peechee.PathReader
    Logger logging.Logger

    config *Config
    muxer *http.ServeMux
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
    logger logging.Logger,
) *http.Server {
    return &http.Server{
        Addr: fmt.Sprintf("%s:%s", host, strconv.Itoa(port)),
        Handler: muxer,
        TLSConfig: tlsConfig,
        TLSNextProto: make(
            map[string]func(*http.Server, *tls.Conn, http.Handler),
        ),
        ErrorLog: logger.GetStdLogger(logging.ERROR),
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
    allRouteControllers := append([]RouteController{rc}, otherRCs)
    methodCallers := MethodsForRouteController(path, allRouteControllers...)
    routeControllerHelper := RouteControllerHelper{
        route: route,
        methodCallers: methodCallers,
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
    s.Logger.Info("Stopping server...").Send()
    err := s.internalServer.Close()
    if err != nil {
        return errors.Wrap(err, "Error while stopping server")
    }

    return nil
}

func (s Server) startInternalServer() error {
    s.Logger.
        Info(fmt.Sprintf(
            "Starting server...",
        )).
        With("host", s.config.Host).
        With("port", s.config.Port).
        With("using_encryption", s.encryptionEnabled).
        Send()
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
        s.Logger.Info("Server shutdown.").Send()
        outCh <- ServerChannelResponse{
            Type: ServerShutdownChannelResponse,
            Error: nil,
        }
    }
}

func (s *Server) startUp() error {
    err := s.startInternalServer()
    if err != http.ErrServerClosed {
        s.Logger.
            Error(fmt.Sprintf("%s", err)).
            With("config", s.internalServer.TLSConfig).
            Send()
        return errors.Wrap(err, "Error while running server")
    }

    s.Logger.Info("Server shutdown.").Send()
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
            ).Send()

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

    pathVariables := rch.route.PathParams(r.URL.Path)
    transactor, err := NewTransactor(
        r, w, pathVariables, self.config, self.Logger, self.defaultEncoding,
    )
    if err != nil {
        self.Logger.Error("Error while creating Transactor.").With(
            "error", err,
        ).Send()
        return responses.ErrorResponse(err)
    }

    for _, middleware := range self.preActionMiddleware {
        data, err := middleware(r.Context(), transactor)
        if err != nil {
            self.Logger.Error("Error in pre-action middleware.").With(
                "error", err,
            ).Send()
            return responses.ErrorResponse(err)
        }
        if data != nil {
            // If we have data from our middleware return early with it.
            return *data
        }
    }

    self.Logger.Info("Handling HTTP request.").
        With("path", r.URL.Path).
        With("method", r.Method).
        With("sequence_id", transactor.SequenceId()).
        Send()

    response := rcc(r.Context(), transactor)

    for _, middleware := range self.postActionMiddleware {
        data, err := middleware(r.Context(), transactor, response)
        if err != nil {
            self.Logger.Error("Error in pre-action middleware.").With(
                "error", err,
            ).Send()
            return responses.ErrorResponse(err)
        }
        if data != nil {
            // If we have data to return early with it.
            return *data
        }
    }

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
                errString := fmt.Sprintf("%+v", rawErr)
                if err, ok := rawErr.(error); ok {
                    wrappedErr := errors.WithStack(err)
                    errString = fmt.Sprintf("%+v", wrappedErr)
                }
                server.Logger.Error("Panic while handling controller.").With(
                    "error", errString,
                ).Send()
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "Internal Server Error")
            }
        }()

        ctx, _ = handleSequenceId(r)
        r = r.WithContext(ctx)

        responseData := handlerFunc(w, r)
        if unexpectedErr := responseData.Error(); unexpectedErr != nil {
            if unexpectedErr.Error() != "" {
                server.Logger.Error(
                    "Unexpected error in controller route.",
                ).With(
                    "error", fmt.Sprint(unexpectedErr),
                ).And("sequence_id", sequenceId).Send()
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
        sequenceId, err := ContextSequenceId(r.Context())
        if err != nil {
            s.Logger.Warn("Unable to extract SequenceId from context.").With(
                "error", err,
            ).Send()
        }
        s.Logger.Info("HTTP request made.").
            With("path", r.URL.Path).
            With("requestor", r.RemoteAddr).
            With("method", r.Method).
            With("sequence_id", sequenceId).
            Send()

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

func (self *Server) defaultHandler() (
    func(w http.ResponseWriter, r *http.Request),
) {
    defaultHandler := func(
        w http.ResponseWriter, r *http.Request,
    ) responses.Data {
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

    return responseProcessor(defaultHandler, self)
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
        logLevels, err := logging.GetLogLevelsForString("warn")
        if err != nil {
            return nil, errors.Wrap(
                err,
                "Error while trying to get log level for logger",
            )
        }
        logger = logging.GetNewLogger(
            "vial", logging.JSON, logging.Stdout, logLevels,
        )
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
