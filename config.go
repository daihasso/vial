package vial

import (
)

type Config struct {
    Swagger struct {
        Path string
    }
    Host string
    Port int
    Tls struct {
        CertPath,
        KeyPath string
        Enabled bool
    }
    Jwt struct {
        EncryptionKey,
        HmacKey string
    }
}

func newConfig() *Config {
    config := &Config{
        Host: "127.0.0.1",
        Port: 8080,
        Tls: struct{
            CertPath,
            KeyPath string
            Enabled bool
        }{
            CertPath: "",
            KeyPath: "",
            Enabled: false,
        },
    }

    return config
}
