package main

import (
	"log"
	"os"
)

type EnvConfigInterface interface {
	IsDev() bool
	Addr() string
}

type EnvConfig struct {
	isDev bool
	addr  string
}

func NewEnv() *EnvConfig {
	log.Println("Loading environment variables...")

	var addr string
	addr = os.Getenv("DEVELOPMENT_PORT")
	return &EnvConfig{addr: addr}
}

func (e *EnvConfig) IsDev() bool {
	return e.isDev
}

func (e *EnvConfig) Addr() string {
	return e.addr
}
