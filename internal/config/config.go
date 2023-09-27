package config

import "time"

type Config struct {
	ServerAddress         string        `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8081"`
	ServerTimeout         time.Duration `env:"SERVER_TIMEOUT" envDefault:"10s"`
	ServerReadBufferSize  int           `env:"SERVER_READ_BUFFER_SIZE" envDefault:"1024"`
	ServerWriteBufferSize int           `env:"SERVER_WRITE_BUFFER_SIZE" envDefault:"1024"`
}
