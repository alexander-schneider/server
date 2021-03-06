package main

import (
	"github.com/getsentry/sentry-go"
	l "github.com/loeffel-io/logger"
	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
)

const (
	port = "8080"

	certFile = "/home/certs/server.crt"
	keyFile  = "/home/certs/server.key"
)

func main() {
	var (
		err error
		zip = new(archiver.Zip)
	)

	// Setup sentry
	if err = sentry.Init(sentry.ClientOptions{Dsn: os.Getenv("SENTRY")}); err != nil {
		log.Fatal(err)
	}

	// Logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logger := &l.Logger{
		Debug:   true,
		RWMutex: new(sync.RWMutex),
	}

	// Config
	maxSize, err := strconv.ParseInt(getenv("MAX_SIZE", "32"), 10, 64)

	if err != nil {
		logger.Error(err)
	}

	token, err := mustenv("TOKEN")

	if err != nil {
		logger.Error(err)
	}

	// SSL
	ssl := true

	for _, path := range []string{certFile, keyFile} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			ssl = false
			break
		}
	}

	// API
	api := &api{
		ssl:     ssl,
		zip:     zip,
		port:    port,
		mode:    "debug",
		maxSize: maxSize,
		token:   token,
		RWMutex: new(sync.RWMutex),
	}

	if err = api.startServer(); err != nil {
		logger.Error(err)
	}
}
