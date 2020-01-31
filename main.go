package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/sirupsen/logrus"
	"github.com/stevenroose/gonfig"
)

var version = "devel"

var logger *logrus.Logger

func main() {
	logger = logrus.New()

	if err := gonfig.Load(&conf, gonfig.Conf{
		EnvPrefix: "ATTO_",
	}); err != nil {
		logger.Fatalf("could not load config: %s", err)
	}

	if level, err := logrus.ParseLevel(conf.LogLevel); err == nil {
		logger.SetLevel(level)
	} else {
		logger.Fatal(err)
	}
	log.SetOutput(logger.Writer())

	logger.WithFields(logrus.Fields{
		"version": version,
	}).Info("starting atto")

	handler := http.StripPrefix(conf.Prefix, http.FileServer(safeDir(conf.Path)))

	if conf.Compress {
		handler = gziphandler.GzipHandler(handler)
	}

	if conf.Canonical.Host != "" {
		handler = canonicalRedirectMiddleware(handler.ServeHTTP)
	}

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Port),
		Handler:           loggingMiddleware(handler.ServeHTTP),
		ReadHeaderTimeout: time.Duration(*conf.Timeout.ReadHeader),
		ReadTimeout:       time.Duration(*conf.Timeout.ReadHeader), // we do not expect any content upload, so headers are enough
		ErrorLog:          log.New(logger.WithField("source", "http.Server").WriterLevel(logrus.WarnLevel), "", 0),
	}

	server.RegisterOnShutdown(func() {
		logger.WithFields(logrus.Fields{
			"timeout": time.Duration(*conf.Timeout.Shutdown),
		}).Info("shutting down gracefully")
	})

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		logger.WithFields(logrus.Fields{
			"port": conf.Port,
		}).Debug("starting server")

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		handleSignals(&server)
	}()

	wg.Wait()
}

func loggingMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		logger.WithFields(logrus.Fields{
			"host":   r.Host,
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("handled request")
	}
}

func canonicalRedirectMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if hostNoPort, _, err := net.SplitHostPort(r.Host); err == nil {
			host = hostNoPort
		}

		if conf.Canonical.Host != host {
			logger.WithFields(logrus.Fields{
				"canonical": conf.Canonical.Host,
			}).Debug("redirecting to canonical host")

			u := r.URL
			u.Host = conf.Canonical.Host

			http.Redirect(w, r, u.String(), conf.Canonical.StatusCode)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func handleSignals(server *http.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*conf.Timeout.Shutdown))

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			logger.Warn("timeout exceeded while shutting down")
		} else {
			logger.WithError(err).Error("error while shutting down server")
		}
	}
}
