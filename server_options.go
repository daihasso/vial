package vial

import (
    "io"

    "github.com/daihasso/slogging"
    "github.com/pkg/errors"
    "daihasso.net/library/tote"
    "github.com/daihasso/peechee"

    "daihasso.net/library/vial/responses"
)

type serverModifier func(*Server) error

type serverOptions struct {
    preActionMiddleware []PreMiddleWare
    postActionMiddleware []PostMiddleWare
    config *Config
    logger logging.Logger
    pathReader *peechee.PathReader
    defaultEncoding responses.EncodingType

    tlsCertData,
    tlsKeyData io.Reader
    useEncryption bool
    serverMods []serverModifier
}

func newServerOptions() *serverOptions {
    return &serverOptions{
        // Default to filesystem only.
        pathReader: peechee.NewPathReader(peechee.WithFilesystem()),
    }
}

// ServerOption is a n option applied to the server.
type ServerOption func(*serverOptions) error

// AddPreActionMiddleware runs the provided middleware(s) against the server
// before every request is handled.
func AddPreActionMiddleware(middlewares ...PreMiddleWare) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.preActionMiddleware = append(
            svOpts.preActionMiddleware, middlewares...,
        )

        return nil
    }
}

// AddPostActionMiddleware runs the provided middleware(s) against the server
// after every request is handled but before it is returned to the requestor.
func AddPostActionMiddleware(middlewares ...PostMiddleWare) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.postActionMiddleware = append(
            svOpts.postActionMiddleware, middlewares...,
        )

        return nil
    }
}

// AddConfig sets the server config to the provided config.
func AddConfig(config *Config) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.config = config

        return nil
    }
}

// AddConfig reads a config from a file under the `vial:` key in a config
// file at the provided path.
func AddConfigFromFile(configFilePath string) ServerOption {
    return func(svOpts *serverOptions) error {
        newConfig := newConfig()
        wrap := struct{Vial *Config}{newConfig}
        err := tote.ReadConfig(&wrap, tote.AddPaths(configFilePath))
        if err != nil {
            return errors.WithStack(err)
        }
        svOpts.config = newConfig

        return nil
    }
}

// AddCustomLogger will set the server logger to the provided logger.
func AddCustomLogger(logger logging.Logger) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.logger = logger

        return nil
    }
}

// AddEncryption adds encryption to the server with the provided cert/key data.
func AddEncryption(certData, keyData io.Reader) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.useEncryption = true
        svOpts.tlsCertData = certData
        svOpts.tlsKeyData = keyData

        return nil
    }
}

// AddEncryptionFilePaths adds encryption to the server with the key/cert data
// at the provided paths.
func AddEncryptionFilePaths(certPath, keyPath string) ServerOption {
    return func(svOpts *serverOptions) error {
        certData, err := svOpts.pathReader.Read(certPath)
        if err != nil {
            return errors.Wrapf(
                err, "Couldn't get certificate data at path '%s'", certPath,
            )
        }
        keyData, err := svOpts.pathReader.Read(keyPath)
        if err != nil {
            return errors.Wrapf(
                err, "Couldn't get private key data at path '%s'", keyPath,
            )
        }
        svOpts.useEncryption = true
        svOpts.tlsCertData = certData
        svOpts.tlsKeyData = keyData

        return nil
    }
}

// AddDefaultSwaggerRoute adds the default swagger route to the server.
// This route is included in NewServerDefault automatically.
func AddDefaultSwaggerRoute(formats ...SwaggerFormat) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.serverMods = append(
            svOpts.serverMods, func(server *Server) error {
                for _, format := range formats {
                    if format == SwaggerYamlFormat {
                        swaggerYamlController := &defaultSwaggerController{
                            false, server.config,
                        }
                        err := server.AddController(
                            "/swagger.yaml", swaggerYamlController,
                        )
                        if err != nil {
                            return errors.Wrap(
                                err,
                                "Error while adding swagger YAML controller",
                            )
                        }
                    }
                    if format == SwaggerJsonFormat {
                        swaggerJsonController := &defaultSwaggerController{
                            true, server.config,
                        }
                        err := server.AddController(
                            "/swagger.json", swaggerJsonController,
                        )
                        if err != nil {
                            return errors.Wrap(
                                err,
                                "Error while adding swagger JSON controller",
                            )
                        }
                    }
                }

                return nil
            },
        )

        return nil
    }
}

// AddDefaultHealthRoute adds an ultra-simple health route that simple returns
// a JSON payload with `{ "healthy": true }`
// This route is included in NewServerDefault automatically.
func AddDefaultHealthRoute() ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.serverMods = append(
            svOpts.serverMods, func(server *Server) error {
                err := server.AddController(
                    "/health", &defaultHealthController{},
                )

                if err != nil {
                    return errors.Wrap(
                        err,
                        "Error while adding default health controller",
                    )
                }

                return nil
            },
        )

        return nil
    }
}

// SetDefaultEncoding will set the server-wide default encoding scheme.
func SetDefaultEncoding(encodingType responses.EncodingType) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.defaultEncoding = encodingType

        return nil
    }
}

// SetPathReader overrides the default PathReader with a custom one.
func SetPathReader(pathReader *peechee.PathReader) ServerOption {
    return func(svOpts *serverOptions) error {
        svOpts.pathReader = pathReader

        return nil
    }
}
